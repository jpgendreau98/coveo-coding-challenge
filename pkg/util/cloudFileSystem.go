package util

import (
	"math"
	"time"
)

type CloudFilesystem interface {
	SetName(value string)
	SetCreationDate(value time.Time)
	SetNbOfFiles(value int64)
	SetSizeOfBucket(value float64)
	SetLastUpdateDate(value time.Time)
	SetCost(value float64)
	SetStorageClass(value StorageClassSizeMap)
	SetRegion(value string)
	ApplySizeConversion(sizeConvrsion float64)
	GetName() string
	GetCreationDate() time.Time
	GetNbOfFiles() int64
	GetSizeOfBucket() float64
	GetLastUpdateDate() time.Time
	GetCost() float64
	GetStorageClass() StorageClassSizeMap
	GetRegion() string
}

type BucketDTO struct {
	Name             string
	CreationDate     time.Time
	NbOfFiles        int64
	SizeOfBucket     float64
	LastUpdateDate   time.Time
	Cost             float64
	StorageClassSize StorageClassSizeMap
	Region           string
}

func NewCloudFileSystem(fsType string) CloudFilesystem {
	switch fsType {
	case "S3":
		return &BucketDTO{}
	default:
		return nil
	}
}

func (bucket *BucketDTO) SetName(value string) {
	bucket.Name = value
}

func (bucket *BucketDTO) SetCreationDate(value time.Time) {
	bucket.CreationDate = value
}

func (bucket *BucketDTO) SetNbOfFiles(value int64) {
	bucket.NbOfFiles = value
}

func (bucket *BucketDTO) SetSizeOfBucket(value float64) {
	bucket.SizeOfBucket = value
}

func (bucket *BucketDTO) SetLastUpdateDate(value time.Time) {
	bucket.LastUpdateDate = value
}

func (bucket *BucketDTO) SetCost(value float64) {
	bucket.Cost = value
}

func (bucket *BucketDTO) SetStorageClass(value StorageClassSizeMap) {
	bucket.StorageClassSize = value
}

func (bucket *BucketDTO) SetRegion(value string) {
	bucket.Region = value
}

func (bucket *BucketDTO) GetName() string {
	return bucket.Name
}

func (bucket *BucketDTO) GetCreationDate() time.Time {
	return bucket.CreationDate
}

func (bucket *BucketDTO) GetNbOfFiles() int64 {
	return bucket.NbOfFiles
}

func (bucket *BucketDTO) GetSizeOfBucket() float64 {
	return bucket.SizeOfBucket
}

func (bucket *BucketDTO) GetLastUpdateDate() time.Time {
	return bucket.LastUpdateDate
}

func (bucket *BucketDTO) GetCost() float64 {
	return bucket.Cost
}

func (bucket *BucketDTO) GetStorageClass() StorageClassSizeMap {
	return bucket.StorageClassSize
}

func (bucket *BucketDTO) GetRegion() string {
	return bucket.Region
}

func (bucket *BucketDTO) ApplySizeConversion(sizeConversion float64) {
	bucket.SizeOfBucket = bucket.SizeOfBucket / math.Pow(float64(1024), sizeConversion)
	for k, v := range bucket.StorageClassSize {
		bucket.StorageClassSize[k] = v / math.Pow(float64(1024), sizeConversion)
	}
}
