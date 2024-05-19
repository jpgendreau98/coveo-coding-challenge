package aws

import (
	"cmp"
	"math"
	"projet-devops-coveo/pkg/util"
	"slices"
)

func TransformByteToGB(size float64, sizeConversion float64) float64 {
	//Transform back size to bytes with the conversion options
	return (size * (math.Pow(1024, sizeConversion))) / (1024 * 1024 * 1024)
}

func SortListBasedOnRegion(buckets []*util.BucketDTO) {
	slices.SortStableFunc(buckets, func(a, b *util.BucketDTO) int {
		return cmp.Compare(a.Region, b.Region)
	})
}

func GetBucketOfTheRegion(buckets []*util.BucketDTO, region string) (regionBuckets []*util.BucketDTO) {
	for _, bucket := range buckets {
		if bucket.Region == region {
			regionBuckets = append(regionBuckets, bucket)
		}
	}
	return regionBuckets
}

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
