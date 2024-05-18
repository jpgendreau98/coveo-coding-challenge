package pkg

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	session               *s3.S3
	totalStorageClassSize StorageClassSize
}

type Bucket struct {
	Name             string
	CreationDate     time.Time
	NbOfFiles        int64
	SizeOfBucket     int64
	LastUpdateDate   string
	Cost             float64
	StorageClassSize StorageClassSize
}

type StorageClassSize map[string]int64

func InitConnection() (*S3, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	})
	if err != nil {
		return nil, err
	}
	return &S3{
		session:               s3.New(session),
		totalStorageClassSize: make(StorageClassSize),
	}, nil
}

func (fs *S3) ListDirectories() []*s3.Bucket {
	output, err := fs.session.ListBuckets(nil)
	if err != nil {
		fmt.Println(err)
	}
	return output.Buckets
}

func (fs *S3) GetObjectMetadata(directories []*s3.Bucket, priceList MasterPriceList) (buckets []*Bucket) {
	for _, bucket := range directories {
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(string(*bucket.Name)),
		}
		storageClassSize := make(StorageClassSize)
		var totalSize int64
		var nbOfFiles int64
		err := fs.session.ListObjectsV2Pages(input,
			func(page *s3.ListObjectsV2Output, lastPage bool) bool {
				nbOfFiles = aws.Int64Value(page.KeyCount)
				for _, obj := range page.Contents {
					totalSize += *obj.Size
					storageClassSize[aws.StringValue(obj.StorageClass)] += *obj.Size
					fs.totalStorageClassSize[aws.StringValue(obj.StorageClass)] += *obj.Size
				}
				return !lastPage
			})
		if err != nil {
			fmt.Println("Error listing objects:", err)
			continue
		}

		outputBucket := &Bucket{
			Name:             string(*bucket.Name),
			CreationDate:     *bucket.CreationDate,
			NbOfFiles:        nbOfFiles,
			SizeOfBucket:     totalSize,
			StorageClassSize: storageClassSize,
		}
		buckets = append(buckets, outputBucket)
	}
	return buckets
}

func (fs *S3) GetBucketPrices(buckets []*Bucket, priceList MasterPriceList) {
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

func ShowCostInString(cost float64) string {
	return fmt.Sprintf("%.16f", cost)
}
