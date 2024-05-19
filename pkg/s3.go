package pkg

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	session               *s3.S3
	totalStorageClassSize StorageClassSize
	region                string
}

type BucketDTO struct {
	Name             string
	CreationDate     time.Time
	NbOfFiles        int64
	SizeOfBucket     int64
	LastUpdateDate   time.Time
	Cost             float64
	StorageClassSize StorageClassSize
	Region           string
}

type StorageClassSize map[string]int64

func InitConnection(region string) (*S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return &S3{
		session:               s3.New(session),
		totalStorageClassSize: make(StorageClassSize),
		region:                region,
	}, nil
}

func (fs *S3) ListDirectories(options CliOptions) (bucketList []*BucketDTO) {
	output, err := fs.session.ListBuckets(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, bucket := range output.Buckets {
		location := fs.GetBucketLocation(bucket)
		if slices.Contains(options.Regions, location) {
			// Filter by names if there's a filter activated
			if options.FilterByName != nil && len(options.FilterByName) != 0 && !slices.Contains(options.FilterByName, aws.StringValue(bucket.Name)) {
				continue
			}

			bucketList = append(bucketList, &BucketDTO{
				Name:         aws.StringValue(bucket.Name),
				Region:       location,
				CreationDate: *bucket.CreationDate,
			})
		}
	}
	return bucketList
}

func (fs *S3) GetBucketLocation(bucket *s3.Bucket) string {
	result, err := fs.session.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		fmt.Println(err)
	}
	if aws.StringValue(result.LocationConstraint) != "" {
		return aws.StringValue(result.LocationConstraint)
	}
	return "us-east-1"
}

func (fs *S3) GetObject(directories []*BucketDTO, priceList MasterPriceList, options CliOptions) (buckets []*BucketDTO) {
	bucketChan := make(chan *BucketDTO, len(directories))
	wg := new(sync.WaitGroup)
	for _, bucket := range directories {
		wg.Add(1)
		go fs.FetchBucket(bucket, bucketChan, options, wg)
	}
	wg.Wait()
	close(bucketChan)
	for bucket := range bucketChan {

		buckets = append(buckets, bucket)
	}
	return buckets
}

func (fs *S3) FetchBucket(bucket *BucketDTO, bucketChan chan (*BucketDTO), options CliOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	if bucket.Region != fs.region {
		return
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket.Name),
	}
	storageClassSize := make(StorageClassSize)
	var totalSize int64
	var nbOfFiles int64
	var lastModifiedBucket time.Time
	loc, _ := time.LoadLocation("Local")

	err := fs.session.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			nbOfFiles = aws.Int64Value(page.KeyCount)
			for _, obj := range page.Contents {
				if options.FilterByStorageClass != nil {
					if !slices.Contains(options.FilterByStorageClass, aws.StringValue(obj.StorageClass)) {
						continue
					}
				}
				totalSize += *obj.Size
				storageClassSize[aws.StringValue(obj.StorageClass)] += *obj.Size
				fs.totalStorageClassSize[aws.StringValue(obj.StorageClass)] += *obj.Size

				if lastModifiedBucket.Before(*obj.LastModified) {
					currentTime := obj.LastModified.In(loc)
					lastModifiedBucket = currentTime
				}
			}
			return !lastPage
		})

	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}
	if !options.ReturnEmptyBuckets {
		if nbOfFiles == 0 && totalSize == 0 {
			return
		}
	}

	bucket.NbOfFiles = nbOfFiles
	bucket.SizeOfBucket = totalSize
	bucket.StorageClassSize = storageClassSize
	bucket.LastUpdateDate = lastModifiedBucket

	bucketChan <- bucket
}

func (fs *S3) SetBucketPrices(buckets []*BucketDTO, priceList MasterPriceList) {
	tierListPrice := GetTierPriceList(fs.totalStorageClassSize, priceList)
	for _, bucket := range buckets {
		var total float64
		for k, v := range bucket.StorageClassSize {
			totalSize := fs.totalStorageClassSize[k]
			price := tierListPrice[k]
			total += (TransformByteToGB(v) / TransformByteToGB(totalSize)) * price
		}
		bucket.Cost = total
	}

}
