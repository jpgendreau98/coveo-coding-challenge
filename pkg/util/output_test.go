package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByRegion(t *testing.T) {
	tests := []struct {
		input  []*BucketDTO
		output map[string][]*BucketDTO
	}{
		{
			input: []*BucketDTO{
				{
					Name:   "test1",
					Region: "ca-central-1",
				}, {
					Name:   "test2",
					Region: "ca-central-1",
				}, {
					Name:   "test3",
					Region: "us-east-1",
				}, {
					Name:   "test4",
					Region: "us-west-2",
				},
			},
			output: map[string][]*BucketDTO{
				"ca-central-1": {
					{
						Name:   "test1",
						Region: "ca-central-1",
					}, {
						Name:   "test2",
						Region: "ca-central-1",
					},
				},
				"us-east-1": {
					{
						Name:   "test3",
						Region: "us-east-1",
					},
				},
				"us-west-2": {
					{
						Name:   "test4",
						Region: "us-west-2",
					},
				},
			},
		},
	}
	for _, test := range tests {
		output := groupByRegion(test.input)
		assert.Equal(t, test.output, output)
	}
}

func TestOrderByInc(t *testing.T) {
	input := []*BucketDTO{
		{
			Name:         "test1",
			Region:       "ca-central-1",
			SizeOfBucket: 934,
			Cost:         123,
		}, {
			Name:         "test2",
			Region:       "ca-central-1",
			SizeOfBucket: 437589374,
			Cost:         1234,
		}, {
			Name:         "test3",
			Region:       "us-east-1",
			SizeOfBucket: 84834,
			Cost:         12345,
		}, {
			Name:         "test4",
			Region:       "us-west-2",
			SizeOfBucket: 934223,
			Cost:         123456,
		},
	}
	tests := []struct {
		input  []*BucketDTO
		output []*BucketDTO
		key    string
	}{
		{
			input: input,
			output: []*BucketDTO{
				{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, {
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
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
			output: []*BucketDTO{
				{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, {
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
			output: []*BucketDTO{
				{
					Name:         "test1",
					Region:       "ca-central-1",
					SizeOfBucket: 934,
					Cost:         123,
				}, {
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
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
	input := []*BucketDTO{
		{
			Name:         "test1",
			Region:       "ca-central-1",
			SizeOfBucket: 934,
			Cost:         123,
		}, {
			Name:         "test2",
			Region:       "ca-central-1",
			SizeOfBucket: 437589374,
			Cost:         1234,
		}, {
			Name:         "test3",
			Region:       "us-east-1",
			SizeOfBucket: 84834,
			Cost:         12345,
		}, {
			Name:         "test4",
			Region:       "us-west-2",
			SizeOfBucket: 934223,
			Cost:         123456,
		},
	}
	tests := []struct {
		input  []*BucketDTO
		output []*BucketDTO
		key    string
	}{
		{
			input: input,
			output: []*BucketDTO{
				{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, {
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
			output: []*BucketDTO{
				{
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, {
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
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
			output: []*BucketDTO{
				{
					Name:         "test4",
					Region:       "us-west-2",
					SizeOfBucket: 934223,
					Cost:         123456,
				}, {
					Name:         "test3",
					Region:       "us-east-1",
					SizeOfBucket: 84834,
					Cost:         12345,
				}, {
					Name:         "test2",
					Region:       "ca-central-1",
					SizeOfBucket: 437589374,
					Cost:         1234,
				}, {
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
