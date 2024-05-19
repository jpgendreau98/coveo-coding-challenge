package aws

import (
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"projet-devops-coveo/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/ratelimit"
)

type S3 struct {
	session               *s3.S3
	totalStorageClassSize *util.StorageClassSize
	region                string
	options               util.CliOptions
	limiter               ratelimit.Limiter
}

func InitConnection(region string, options util.CliOptions, globalStorageClass *util.StorageClassSize, limiter ratelimit.Limiter) (*S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return &S3{
		session:               s3.New(session),
		totalStorageClassSize: globalStorageClass,
		region:                region,
		options:               options,
		limiter:               limiter,
	}, nil
}

func (fs *S3) GetBucketsFiltered() (bucketList []*util.BucketDTO) {
	fs.limiter.Take()
	output, err := fs.session.ListBuckets(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	return fs.FilterBuckets(output.Buckets)
}

func (fs *S3) FilterBuckets(buckets []*s3.Bucket) (bucketList []*util.BucketDTO) {
	for _, bucket := range buckets {
		location := fs.GetBucketLocation(bucket)
		if slices.Contains(fs.options.Regions, location) {
			// Filter by names if there's a filter activated
			if fs.options.FilterByName != nil && len(fs.options.FilterByName) != 0 && !slices.Contains(fs.options.FilterByName, aws.StringValue(bucket.Name)) {
				continue
			}

			bucketList = append(bucketList, &util.BucketDTO{
				Name:         aws.StringValue(bucket.Name),
				Region:       location,
				CreationDate: *bucket.CreationDate,
			})
		}
	}
	return bucketList
}

func (fs *S3) ListObjectsInBucket(regionBucket []*util.BucketDTO, region string, priceList MasterPriceList, wg *sync.WaitGroup, bucketChan chan ([]*util.BucketDTO)) {
	defer wg.Done()
	DTOBuckets := fs.GetObject(regionBucket, priceList)
	bucketChan <- DTOBuckets
}

func (fs *S3) GetBucketLocation(bucket *s3.Bucket) string {
	fs.limiter.Take()
	result, err := fs.session.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		fmt.Println(err)
	}
	if aws.StringValue(result.LocationConstraint) != "" {
		return aws.StringValue(result.LocationConstraint)
	}
	return "us-east-1" //Problem with the API response that return empty strings when region is us-east-1
}

func (fs *S3) GetObject(directories []*util.BucketDTO, priceList MasterPriceList) (buckets []*util.BucketDTO) {
	bucketChan := make(chan *util.BucketDTO, len(directories))
	concurrencyThrottle := make(chan int, fs.options.Threading)
	wg := new(sync.WaitGroup)
	for _, bucket := range directories {
		wg.Add(1)
		concurrencyThrottle <- 1
		go fs.FetchBucket(bucket, bucketChan, wg, concurrencyThrottle)
	}
	wg.Wait()
	close(bucketChan)
	for bucket := range bucketChan {

		buckets = append(buckets, bucket)
	}
	return buckets
}

func (fs *S3) FetchBucket(bucket *util.BucketDTO, bucketChan chan (*util.BucketDTO), wg *sync.WaitGroup, concurrencyThrottle chan (int)) {
	defer wg.Done()
	defer func() {
		<-concurrencyThrottle
	}()
	if bucket.Region != fs.region {
		return
	}

	storageClassSize := make(util.StorageClassSizeMap)
	var totalSize int64
	var nbOfFiles int64
	var lastModifiedBucket time.Time
	loc, _ := time.LoadLocation("Local")

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket.Name),
	}

	err := fs.session.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			nbOfFiles = aws.Int64Value(page.KeyCount)
			for _, obj := range page.Contents {
				if fs.options.FilterByStorageClass != nil {
					if len(fs.options.FilterByStorageClass) != 0 && fs.options.FilterByStorageClass != nil &&
						!slices.Contains(fs.options.FilterByStorageClass, aws.StringValue(obj.StorageClass)) {
						continue
					}
				}

				totalSize += *obj.Size
				storageClassSize[aws.StringValue(obj.StorageClass)] += float64(*obj.Size)
				fs.totalStorageClassSize.Mutex.Lock()
				fs.totalStorageClassSize.SizeMap[aws.StringValue(obj.StorageClass)] += float64(*obj.Size)
				fs.totalStorageClassSize.Mutex.Unlock()

				if lastModifiedBucket.Before(*obj.LastModified) {
					currentTime := obj.LastModified.In(loc)
					lastModifiedBucket = currentTime
				}
			}
			fs.limiter.Take()
			return !lastPage
		})

	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}
	if !fs.options.ReturnEmptyBuckets {
		if nbOfFiles == 0 && totalSize == 0 {
			return
		}
	}

	bucket.NbOfFiles = nbOfFiles
	bucket.SizeOfBucket = float64(totalSize) / (math.Pow(float64(1024), fs.options.OutputOptions.SizeConversion))
	bucket.StorageClassSize.SizeMap = storageClassSize
	bucket.LastUpdateDate = lastModifiedBucket

	bucketChan <- bucket
}

func (fs *S3) SetBucketPrices(buckets []*util.BucketDTO, priceList MasterPriceList) {
	tierListPrice := GetTierPriceList(fs.totalStorageClassSize.SizeMap, priceList)
	for _, bucket := range buckets {
		var total float64
		for k, v := range bucket.StorageClassSize.SizeMap {
			totalSize := fs.totalStorageClassSize.SizeMap[k]
			price := tierListPrice[k]
			total += (TransformByteToGB(v, fs.options.OutputOptions.SizeConversion) / TransformByteToGB(totalSize, fs.options.OutputOptions.SizeConversion)) * price
		}
		bucket.Cost = total
	}
}
