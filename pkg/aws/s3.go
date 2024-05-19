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

// Establish connection with S3 services
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

// List All buckets and  returns a filtered list based on filters (name, region)
func (fs *S3) GetBucketsFiltered() (bucketList []*util.BucketDTO) {
	fs.limiter.Take()
	output, err := fs.session.ListBuckets(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	return fs.FilterBuckets(output.Buckets)
}

// Fetch location of bucket and filter if not in wanted region or not included by name
func (fs *S3) FilterBuckets(buckets []*s3.Bucket) (bucketList []*util.BucketDTO) {
	for _, bucket := range buckets {
		location := fs.GetBucketLocation(bucket)
		bucket := fs.filterbucket(location, aws.StringValue(bucket.Name), *bucket.CreationDate)
		if bucket != nil {
			bucketList = append(bucketList, bucket)
		}
	}
	return bucketList
}

// Filter Bucket For region and name
func (fs *S3) filterbucket(location string, bucketName string, bucketCreationDate time.Time) *util.BucketDTO {
	if slices.Contains(fs.options.Regions, location) {
		// Filter by names if there's a filter activated
		if fs.options.FilterByName != nil && len(fs.options.FilterByName) != 0 && !slices.Contains(fs.options.FilterByName, bucketName) {
			return nil
		}
		return &util.BucketDTO{
			Name:         bucketName,
			Region:       location,
			CreationDate: bucketCreationDate,
		}
	}
	return nil
}

// List all objects in a specific region
func (fs *S3) ListObjectsInBucket(regionBucket []*util.BucketDTO, region string, priceList MasterPriceList, wg *sync.WaitGroup, bucketChan chan ([]*util.BucketDTO)) {
	defer wg.Done()
	DTOBuckets := fs.GetObject(regionBucket, priceList)
	bucketChan <- DTOBuckets
}

// Get bucket Location
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

// Get objects of a buckets
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
		//Filter empty bucket if flag is false
		if fs.options.OmitEmpty && bucket.NbOfFiles == 0 {
			continue
		}
		buckets = append(buckets, bucket)
	}
	return buckets
}

// Get all object of a bucket
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
	//Recursively, list objects in a bucket and build the bucket metadata at the same time.
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
				storageClassSize[aws.StringValue(obj.StorageClass)] += (float64(*obj.Size))

				fs.totalStorageClassSize.Mutex.Lock()
				fs.totalStorageClassSize.SizeMap[fs.region][aws.StringValue(obj.StorageClass)] += float64(*obj.Size)
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

	bucket.NbOfFiles = nbOfFiles
	bucket.SizeOfBucket = float64(totalSize) / (math.Pow(float64(1024), fs.options.OutputOptions.SizeConversion))
	bucket.StorageClassSize = storageClassSize
	bucket.LastUpdateDate = lastModifiedBucket

	bucketChan <- bucket
}

// Set the bucket cost based on total cost of S3 Service
func (fs *S3) SetBucketCost(buckets []*util.BucketDTO, priceList MasterPriceList) {
	tierListPrice := GetTierPriceList(fs.totalStorageClassSize.SizeMap[fs.region], priceList[fs.region], fs.options.OutputOptions.SizeConversion)
	for _, bucket := range buckets {
		var total float64
		for k, v := range bucket.StorageClassSize {
			totalSize := float64(fs.totalStorageClassSize.SizeMap[fs.region][k]) / (math.Pow(float64(1024), fs.options.OutputOptions.SizeConversion))
			total += (TransformSizeToGB(v, fs.options.OutputOptions.SizeConversion) / TransformSizeToGB(totalSize, fs.options.OutputOptions.SizeConversion)) *
				(tierListPrice[k] * TransformSizeToGB(totalSize, fs.options.OutputOptions.SizeConversion))
		}
		bucket.Cost = total
	}
}
