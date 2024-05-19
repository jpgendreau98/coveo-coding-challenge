package pkg

import (
	"cmp"
	"slices"
)

type CliOptions struct {
	FilterByName         []string
	FilterByStorageClass []string
	ReturnEmptyBuckets   bool
	Regions              []string
}

func TransformByteToGB(size int64) float64 {
	return float64(size) / (1024 * 1024 * 1024)
}

func SortListBasedOnRegion(buckets []*BucketDTO) {
	slices.SortStableFunc(buckets, func(a, b *BucketDTO) int {
		return cmp.Compare(a.Region, b.Region)
	})
}

func GetBucketOfTheRegion(buckets []*BucketDTO, region string) (regionBuckets []*BucketDTO) {
	for _, bucket := range buckets {
		if bucket.Region == region {
			regionBuckets = append(regionBuckets, bucket)
		}
	}
	return regionBuckets
}

func RemoveScrappedBucketFromList(scrappedBuckets []*BucketDTO, bucketList []*BucketDTO) []*BucketDTO {
	for _, scrascrappedBucket := range scrappedBuckets {
		i := slices.IndexFunc(bucketList, func(a *BucketDTO) bool {
			return a.Name == scrascrappedBucket.Name
		})
		if i >= 0 {
			bucketList = slices.Delete(bucketList, i, i+1)
		}
	}
	return bucketList
}
