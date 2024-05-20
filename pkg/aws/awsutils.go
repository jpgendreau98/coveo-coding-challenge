package aws

import (
	"cmp"
	"projet-devops-coveo/pkg/util"
	"slices"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Transform from size to GB based on the conversion from user
func TransformSizeToGB(size float64) float64 {
	return (size) / (1024 * 1024 * 1024)
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

func GetStorageClassConstant(value s3types.ObjectStorageClass) string {
	switch value {
	case s3types.ObjectStorageClassStandard:
		return "STANDARD"
	case s3types.ObjectStorageClassReducedRedundancy:
		return "REDUCED_REDUNDANCY"
	case s3types.ObjectStorageClassGlacier:
		return "GLACIER"
	case s3types.ObjectStorageClassStandardIa:
		return "STANDARD_IA"
	case s3types.ObjectStorageClassOnezoneIa:
		return "ONEZONE_IA"
	case s3types.ObjectStorageClassIntelligentTiering:
		return "INTELLIGENT_TIERING"
	case s3types.ObjectStorageClassDeepArchive:
		return "DEEP_ARCHIVE"
	case s3types.ObjectStorageClassOutposts:
		return "OUTPOSTS"
	case s3types.ObjectStorageClassGlacierIr:
		return "GLACIER_IR"
	case s3types.ObjectStorageClassSnow:
		return "SNOW"
	case s3types.ObjectStorageClassExpressOnezone:
		return "EXPRESS_ONEZONE"
	default:
		return ""
	}
}

func GetBucketLocationConstant(value s3types.BucketLocationConstraint) string {
	switch value {
	case s3types.BucketLocationConstraintAfSouth1:
		return "af-south-1"
	case s3types.BucketLocationConstraintApEast1:
		return "ap-east-1"
	case s3types.BucketLocationConstraintApNortheast1:
		return "ap-northeast-1"
	case s3types.BucketLocationConstraintApNortheast2:
		return "ap-northeast-2"
	case s3types.BucketLocationConstraintApNortheast3:
		return "ap-northeast-3"
	case s3types.BucketLocationConstraintApSouth1:
		return "ap-south-1"
	case s3types.BucketLocationConstraintApSouth2:
		return "ap-south-2"
	case s3types.BucketLocationConstraintApSoutheast1:
		return "ap-southeast-1"
	case s3types.BucketLocationConstraintApSoutheast2:
		return "ap-southeast-2"
	case s3types.BucketLocationConstraintApSoutheast3:
		return "ap-southeast-3"
	case s3types.BucketLocationConstraintCaCentral1:
		return "ca-central-1"
	case s3types.BucketLocationConstraintCnNorth1:
		return "cn-north-1"
	case s3types.BucketLocationConstraintCnNorthwest1:
		return "cn-northwest-1"
	case s3types.BucketLocationConstraintEu:
		return "EU"
	case s3types.BucketLocationConstraintEuCentral1:
		return "eu-central-1"
	case s3types.BucketLocationConstraintEuNorth1:
		return "eu-north-1"
	case s3types.BucketLocationConstraintEuSouth1:
		return "eu-south-1"
	case s3types.BucketLocationConstraintEuSouth2:
		return "eu-south-2"
	case s3types.BucketLocationConstraintEuWest1:
		return "eu-west-1"
	case s3types.BucketLocationConstraintEuWest2:
		return "eu-west-2"
	case s3types.BucketLocationConstraintEuWest3:
		return "eu-west-3"
	case s3types.BucketLocationConstraintMeSouth1:
		return "me-south-1"
	case s3types.BucketLocationConstraintSaEast1:
		return "sa-east-1"
	case s3types.BucketLocationConstraintUsEast2:
		return "us-east-2"
	case s3types.BucketLocationConstraintUsGovEast1:
		return "s-gov-east-1"
	case s3types.BucketLocationConstraintUsGovWest1:
		return "s-gov-west-1"
	case s3types.BucketLocationConstraintUsWest1:
		return "us-west-1"
	case s3types.BucketLocationConstraintUsWest2:
		return "us-west-2"
	default:
		return "us-east-1"
	}
}
