package aws

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"projet-devops-coveo/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/ratelimit"
)

type S3 struct {
	session               *AwsClient
	totalStorageClassSize *util.StorageClassSize
	region                string
	options               util.CliOptions
	limiter               ratelimit.Limiter
}

// Establish connection with S3 services
func InitConnection(region string, options util.CliOptions, globalStorageClass *util.StorageClassSize, limiter ratelimit.Limiter) (*S3, error) {
	awsClient, err := NewAwsClient(region, limiter)
	if err != nil {
		return nil, err
	}

	return &S3{
		session:               awsClient,
		totalStorageClassSize: globalStorageClass,
		region:                region,
		options:               options,
		limiter:               limiter,
	}, nil
}

// List All buckets and  returns a filtered list based on filters (name, region)
func (fs *S3) GetBucketsFiltered() (bucketList []*util.BucketDTO) {
	output, err := fs.session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err.Error())
	}

	return fs.FilterBuckets(output.Buckets)
}

// Fetch location of bucket and filter if not in wanted region or not included by name
func (fs *S3) FilterBuckets(buckets []types.Bucket) (bucketList []*util.BucketDTO) {
	for _, bucket := range buckets {
		location, err := fs.session.GetBucketLocation(&s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			fmt.Println(err)
		}
		bucket := fs.filterbucket(location, *bucket.Name, *bucket.CreationDate)
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
	//Recursively, list objects in a bucket and build the bucket metadata at the same time.
	ctx := context.Background()
	paginator := s3.NewListObjectsV2Paginator(fs.session.s3, &s3.ListObjectsV2Input{
		Bucket: &bucket.Name,
	})

	// Iterate through pages
	for paginator.HasMorePages() {
		// Get the next page of objects
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			fmt.Println(err)
		}

		// Process objects
		for _, obj := range resp.Contents {
			if fs.options.FilterByStorageClass != nil {
				if len(fs.options.FilterByStorageClass) != 0 && fs.options.FilterByStorageClass != nil &&
					!slices.Contains(fs.options.FilterByStorageClass, GetStorageClassConstant(obj.StorageClass)) {
					continue
				}
			}
			nbOfFiles += 1
			totalSize += *obj.Size
			storageClassSize[GetStorageClassConstant(obj.StorageClass)] += (float64(*obj.Size))

			fs.totalStorageClassSize.Mutex.Lock()
			fs.totalStorageClassSize.SizeMap[fs.region][GetStorageClassConstant(obj.StorageClass)] += float64(*obj.Size)
			fs.totalStorageClassSize.Mutex.Unlock()

			if lastModifiedBucket.Before(*obj.LastModified) {
				currentTime := obj.LastModified.In(loc)
				lastModifiedBucket = currentTime
			}
		}
	}

	bucket.NbOfFiles = nbOfFiles
	bucket.SizeOfBucket = float64(totalSize)
	bucket.StorageClassSize = storageClassSize
	bucket.LastUpdateDate = lastModifiedBucket

	bucketChan <- bucket
}

// Set the bucket cost based on total cost of S3 Service
func (fs *S3) SetBucketCost(buckets []*util.BucketDTO, priceList MasterPriceList) {
	tierListPrice := GetTierPriceList(fs.totalStorageClassSize.SizeMap[fs.region], priceList[fs.region])
	for _, bucket := range buckets {
		var total float64
		for k, v := range bucket.StorageClassSize {
			totalSize := float64(fs.totalStorageClassSize.SizeMap[fs.region][k])
			total += (TransformSizeToGB(v) / TransformSizeToGB(totalSize)) * (tierListPrice[k] * TransformSizeToGB(totalSize))
		}
		bucket.Cost = total
	}
}
