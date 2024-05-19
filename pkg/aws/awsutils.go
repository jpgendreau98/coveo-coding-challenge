package aws

import (
	"cmp"
	"math"
	"projet-devops-coveo/pkg/util"
	"slices"
)

// Transform from size to GB based on the conversion from user
func TransformSizeToGB(size float64, sizeConversion float64) float64 {
	return (size * (math.Pow(1024, sizeConversion))) / (1024 * 1024 * 1024)
}

// Sort list of buckets based on regions
func SortListBasedOnRegion(buckets []*util.BucketDTO) {
	slices.SortStableFunc(buckets, func(a, b *util.BucketDTO) int {
		return cmp.Compare(a.Region, b.Region)
	})
}

// Get all bucket of the specified region.
func GetBucketsOfRegion(buckets []*util.BucketDTO, region string) (regionBuckets []*util.BucketDTO) {
	for _, bucket := range buckets {
		if bucket.Region == region {
			regionBuckets = append(regionBuckets, bucket)
		}
	}
	return regionBuckets
}

// Remove buckets from a list of buckets. Reduce size of array.
// It has to be sorted beforehand so performance is good.
func RemoveScrappedBucketFromList(scrappedBuckets []*util.BucketDTO, bucketList []*util.BucketDTO) []*util.BucketDTO {
	for _, scrascrappedBucket := range scrappedBuckets {
		i := slices.IndexFunc(bucketList, func(a *util.BucketDTO) bool {
			return a.Name == scrascrappedBucket.Name
		})
		if i >= 0 {
			bucketList = slices.Delete(bucketList, i, i+1)
		}
	}
	return bucketList
}
