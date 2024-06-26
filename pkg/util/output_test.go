package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByRegion(t *testing.T) {
	tests := []struct {
		input  []CloudFilesystem
		output map[string][]CloudFilesystem
	}{
		{
			input: []CloudFilesystem{
				&BucketDTO{
					Name:   "test1",
					Region: "ca-central-1",
				}, &BucketDTO{
					Name:   "test2",
					Region: "ca-central-1",
				}, &BucketDTO{
					Name:   "test3",
					Region: "us-east-1",
				}, &BucketDTO{
					Name:   "test4",
					Region: "us-west-2",
				},
			},
			output: map[string][]CloudFilesystem{
				"ca-central-1": {
					&BucketDTO{
						Name:   "test1",
						Region: "ca-central-1",
					}, &BucketDTO{
						Name:   "test2",
						Region: "ca-central-1",
					},
				},
				"us-east-1": {
					&BucketDTO{
						Name:   "test3",
						Region: "us-east-1",
					},
				},
				"us-west-2": {
					&BucketDTO{
						Name:   "test4",
						Region: "us-west-2",
					},
				},
			},
		},
	}
	for _, test := range tests {
		output := applyOutputOptions(test.input, OutputOptions{
			GroupBy: "region",
		})
		assert.Equal(t, test.output, output)
	}
}

func TestOrderByInc(t *testing.T) {
	input := []CloudFilesystem{
		&BucketDTO{
			Name:         "test1",
			Region:       "ca-central-1",
			SizeOfBucket: 934,
			Cost:         123,
		}, &BucketDTO{
			Name:         "test2",
			Region:       "ca-central-1",
			SizeOfBucket: 437589374,
			Cost:         1234,
		}, &BucketDTO{
			Name:         "test3",
			Region:       "us-east-1",
			SizeOfBucket: 84834,
			Cost:         12345,
		}, &BucketDTO{
			Name:         "test4",
			Region:       "us-west-2",
			SizeOfBucket: 934223,
			Cost:         123456,
		},
	}
	tests := []struct {
		input  []CloudFilesystem
		output []CloudFilesystem
		key    string
	}{
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, &BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				},
			},
			key: "name",
		},
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, &BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				},
			},
			key: "size",
		},
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, &BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				},
			},
			key: "cost",
		},
	}
	for _, test := range tests {
		output := orderByInc(test.key, test.input)
		assert.Equal(t, test.output, output)
	}
}

func TestOrderByDec(t *testing.T) {
	input := []CloudFilesystem{
		&BucketDTO{
			Name:         "test1",
			Region:       "ca-central-1",
			SizeOfBucket: 934,
			Cost:         123,
		}, &BucketDTO{
			Name:         "test2",
			Region:       "ca-central-1",
			SizeOfBucket: 437589374,
			Cost:         1234,
		}, &BucketDTO{
			Name:         "test3",
			Region:       "us-east-1",
			SizeOfBucket: 84834,
			Cost:         12345,
		}, &BucketDTO{
			Name:         "test4",
			Region:       "us-west-2",
			SizeOfBucket: 934223,
			Cost:         123456,
		},
	}
	tests := []struct {
		input  []CloudFilesystem
		output []CloudFilesystem
		key    string
	}{
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, &BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				},
			},
			key: "name",
		},
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, &BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				},
			},
			key: "size",
		},
		{
			input: input,
			output: []CloudFilesystem{
				&BucketDTO{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, &BucketDTO{
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, &BucketDTO{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, &BucketDTO{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				},
			},
			key: "cost",
		},
	}
	for _, test := range tests {
		output := orderByDec(test.key, test.input)
		assert.Equal(t, test.output, output)
	}
}
