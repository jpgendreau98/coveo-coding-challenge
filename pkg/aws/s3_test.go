package aws

import (
	"projet-devops-coveo/pkg/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timeMock = time.Now()

func TestFilterbucket(t *testing.T) {
	tests := []struct {
		name               string
		expectedOutput     *util.BucketDTO
		filterByName       []string
		regions            []string
		bucketLocation     string
		bucketName         string
		bucketCreationDate time.Time
	}{
		{
			name:               "Test Filter By Name",
			filterByName:       []string{"poc-1", "poc-2"},
			regions:            []string{"ca-central-1"},
			bucketLocation:     "ca-central-1",
			bucketName:         "poc-3",
			bucketCreationDate: timeMock,
			expectedOutput:     nil,
		},
		{
			name:               "Test Filter By Name",
			filterByName:       []string{"poc-1", "poc-2"},
			regions:            []string{"ca-central-1"},
			bucketLocation:     "ca-central-1",
			bucketName:         "poc-2",
			bucketCreationDate: timeMock,
			expectedOutput: &util.BucketDTO{
				Name:         "poc-2",
				CreationDate: timeMock,
				Region:       "ca-central-1",
			},
		},
		{
			name:               "Test Filter By Region",
			filterByName:       []string{"poc-1", "poc-3"},
			regions:            []string{"ca-central-1"},
			bucketLocation:     "us-east-1",
			bucketName:         "poc-3",
			bucketCreationDate: timeMock,
			expectedOutput:     nil,
		},
		{
			name:               "Test Filter By Region",
			filterByName:       []string{"poc-1", "poc-3"},
			regions:            []string{"ca-central-1", "us-east-1"},
			bucketLocation:     "us-east-1",
			bucketName:         "poc-3",
			bucketCreationDate: timeMock,
			expectedOutput: &util.BucketDTO{
				Name:         "poc-3",
				CreationDate: timeMock,
				Region:       "us-east-1",
			},
		},
	}
	for _, test := range tests {
		fs := &S3{
			options: util.CliOptions{
				Regions:      test.regions,
				FilterByName: test.filterByName,
			},
		}
		output := fs.filterbucket(test.bucketLocation, test.bucketName, test.bucketCreationDate)
		assert.Equal(t, test.expectedOutput, output)
	}
}

func TestSetBucketCost(t *testing.T) {
	tests := []struct {
		name           string
		buckets        []*util.BucketDTO
		priceList      MasterPriceList
		expectedOutput []*util.BucketDTO
	}{
		{
			name: "Small test",
			buckets: []*util.BucketDTO{
				{
					Name:         "Poc-1",
					SizeOfBucket: float64(50000000),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(50000000),
					},
					Region: "ca-central-1",
				},
			},
			priceList: MasterPriceList{
				"ca-central-1": MockProductPriceList,
			},
			expectedOutput: []*util.BucketDTO{
				{
					Name:         "Poc-1",
					SizeOfBucket: float64(50000000),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(50000000),
					},
					Region: "ca-central-1",
					Cost:   0.011641532182693481,
				},
			},
		},
		{
			name: "2 buckets test",
			buckets: []*util.BucketDTO{
				{
					Name:         "Poc-1",
					SizeOfBucket: float64(50000000),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(50000000),
					},
					Region: "ca-central-1",
				},
				{
					Name:         "Poc-2",
					SizeOfBucket: float64(5000033),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(5000033),
					},
					Region: "ca-central-1",
				},
			},
			priceList: MasterPriceList{
				"ca-central-1": MockProductPriceList,
			},
			expectedOutput: []*util.BucketDTO{
				{
					Name:         "Poc-1",
					SizeOfBucket: float64(50000000),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(50000000),
					},
					Region: "ca-central-1",
					Cost:   0.011641532182693481,
				},
				{
					Name:         "Poc-2",
					SizeOfBucket: float64(5000033),
					StorageClassSize: util.StorageClassSizeMap{
						"STANDARD": float64(5000033),
					},
					Region: "ca-central-1",
					Cost:   0.0011641609016805887,
				},
			},
		},
	}
	fs := &S3{
		totalStorageClassSize: &util.StorageClassSize{
			SizeMap: util.RegionsStorageMap{
				"ca-central-1": map[string]float64{
					S3_STORAGE_CLASS_STANDARD: 50000000000,
					S3_STORAGE_CLASS_GLACIER:  2000000,
				},
			},
		},
		options: util.CliOptions{
			OutputOptions: &util.OutputOptions{
				SizeConversion: 0,
			},
		},
		region: "ca-central-1",
	}
	for _, test := range tests {
		fs.SetBucketCost(test.buckets, test.priceList)
		assert.Equal(t, test.expectedOutput, test.buckets)
	}
}
