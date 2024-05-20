package aws

import (
	"projet-devops-coveo/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortListBasedOnRegion(t *testing.T) {
	var expectedInput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "One",
			Region: "ca-central-1",
		}, &util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		}, &util.BucketDTO{
			Name:   "Three",
			Region: "us-east-1",
		}, &util.BucketDTO{
			Name:   "Four",
			Region: "us-west-1",
		}, &util.BucketDTO{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var expectedOutput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "One",
			Region: "ca-central-1",
		}, &util.BucketDTO{
			Name:   "Five",
			Region: "ch-central-2",
		}, &util.BucketDTO{
			Name:   "Three",
			Region: "us-east-1",
		}, &util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		}, &util.BucketDTO{
			Name:   "Four",
			Region: "us-west-1",
		},
	}
	SortListBasedOnRegion(expectedInput)
	assert.Equal(t, expectedInput, expectedOutput)

}

func TestGetBucketsOfRegion(t *testing.T) {
	var expectedInput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "One",
			Region: "ca-central-1",
		}, &util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		}, &util.BucketDTO{
			Name:   "Three",
			Region: "us-east-1",
		}, &util.BucketDTO{
			Name:   "Four",
			Region: "us-west-1",
		}, &util.BucketDTO{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var expectedOutput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		},
	}
	region := "us-east-2"
	output := GetBucketsOfRegion(expectedInput, region)
	assert.Equal(t, expectedOutput, output)
}

func TestRemoveScrappedBucketList(t *testing.T) {
	var expectedInput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "One",
			Region: "ca-central-1",
		}, &util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		}, &util.BucketDTO{
			Name:   "Three",
			Region: "us-east-1",
		}, &util.BucketDTO{
			Name:   "Four",
			Region: "us-west-1",
		}, &util.BucketDTO{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var scappredBucket = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "One",
			Region: "ca-central-1",
		}, &util.BucketDTO{
			Name:   "Three",
			Region: "us-east-1",
		}, &util.BucketDTO{
			Name:   "Four",
			Region: "us-west-1",
		},
	}
	var expectedoutput = []util.CloudFilesystem{
		&util.BucketDTO{
			Name:   "Two",
			Region: "us-east-2",
		},
		&util.BucketDTO{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	output := RemoveScrappedBucketFromList(scappredBucket, expectedInput)
	assert.Equal(t, expectedoutput, output)
}
