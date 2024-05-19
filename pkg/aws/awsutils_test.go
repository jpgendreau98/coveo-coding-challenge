package aws

import (
	"projet-devops-coveo/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortListBasedOnRegion(t *testing.T) {
	var expectedInput = []*util.BucketDTO{
		{
			Name:   "One",
			Region: "ca-central-1",
		}, {
			Name:   "Two",
			Region: "us-east-2",
		}, {
			Name:   "Three",
			Region: "us-east-1",
		}, {
			Name:   "Four",
			Region: "us-west-1",
		},
		{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var expectedOutput = []*util.BucketDTO{
		{
			Name:   "One",
			Region: "ca-central-1",
		}, {
			Name:   "Five",
			Region: "ch-central-2",
		}, {
			Name:   "Three",
			Region: "us-east-1",
		}, {
			Name:   "Two",
			Region: "us-east-2",
		}, {
			Name:   "Four",
			Region: "us-west-1",
		},
	}
	SortListBasedOnRegion(expectedInput)
	assert.Equal(t, expectedInput, expectedOutput)

}

func TestGetBucketOfTheRegion(t *testing.T) {
	var expectedInput = []*util.BucketDTO{
		{
			Name:   "One",
			Region: "ca-central-1",
		}, {
			Name:   "Two",
			Region: "us-east-2",
		}, {
			Name:   "Three",
			Region: "us-east-1",
		}, {
			Name:   "Four",
			Region: "us-west-1",
		},
		{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var expectedOutput = []*util.BucketDTO{
		{
			Name:   "Two",
			Region: "us-east-2",
		},
	}
	region := "us-east-2"
	output := GetBucketOfTheRegion(expectedInput, region)
	assert.Equal(t, expectedOutput, output)
}

func TestRemoveScrappedBucketList(t *testing.T) {
	var expectedInput = []*util.BucketDTO{
		{
			Name:   "One",
			Region: "ca-central-1",
		}, {
			Name:   "Two",
			Region: "us-east-2",
		}, {
			Name:   "Three",
			Region: "us-east-1",
		}, {
			Name:   "Four",
			Region: "us-west-1",
		},
		{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	var scappredBucket = []*util.BucketDTO{
		{
			Name:   "One",
			Region: "ca-central-1",
		}, {
			Name:   "Three",
			Region: "us-east-1",
		}, {
			Name:   "Four",
			Region: "us-west-1",
		},
	}
	var expectedoutput = []*util.BucketDTO{
		{
			Name:   "Two",
			Region: "us-east-2",
		},
		{
			Name:   "Five",
			Region: "ch-central-2",
		},
	}
	output := RemoveScrappedBucketFromList(scappredBucket, expectedInput)
	assert.Equal(t, expectedoutput, output)
}
