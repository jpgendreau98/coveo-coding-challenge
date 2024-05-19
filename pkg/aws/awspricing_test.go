package aws

import (
	"projet-devops-coveo/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

var MockProductPriceList = ProductPriceList{
	"Amazon Glacier": PriceList{
		Terms: struct {
			OnDemand map[string]TermsAttributes "json:\"OnDemand,omitempty\""
		}{
			OnDemand: map[string]TermsAttributes{
				"1234": {
					PriceDimensions: map[string]PriceDimension{
						"5678": {
							BeginRange: "0",
							EndRange:   "400",
							PricePerUnit: struct {
								Usd string "json:\"USD,omitempty\""
							}{
								Usd: "0.25",
							},
						},
						"8901": {
							BeginRange: "400",
							EndRange:   "Inf",
							PricePerUnit: struct {
								Usd string "json:\"USD,omitempty\""
							}{
								Usd: "0.20",
							},
						},
					},
				},
			},
		},
	},
	"Standard": PriceList{
		Terms: struct {
			OnDemand map[string]TermsAttributes "json:\"OnDemand,omitempty\""
		}{
			OnDemand: map[string]TermsAttributes{
				"1234": {
					PriceDimensions: map[string]PriceDimension{
						"5678": {
							BeginRange: "0",
							EndRange:   "51200",
							PricePerUnit: struct {
								Usd string "json:\"USD,omitempty\""
							}{
								Usd: "0.25",
							},
						},
						"8901": {
							BeginRange: "51200",
							EndRange:   "Inf",
							PricePerUnit: struct {
								Usd string "json:\"USD,omitempty\""
							}{
								Usd: "0.20",
							},
						},
					},
				},
			},
		},
	},
}

func TestGetTierPriceList(t *testing.T) {
	tests := []struct {
		name                  string
		totalStorageClassSize util.StorageClassSizeMap
		priceList             ProductPriceList
		expectedOutput        map[string]float64
	}{
		{
			name: "Test First Tier",
			totalStorageClassSize: util.StorageClassSizeMap{
				S3_STORAGE_CLASS_STANDARD: 0.1,
				S3_STORAGE_CLASS_GLACIER:  0.2,
			},
			priceList: MockProductPriceList,
			expectedOutput: map[string]float64{
				S3_STORAGE_CLASS_STANDARD: 0.25,
				S3_STORAGE_CLASS_GLACIER:  0.25,
			},
		},
		{
			name: "Test Second Tier",
			totalStorageClassSize: util.StorageClassSizeMap{
				S3_STORAGE_CLASS_STANDARD: 0.1,
				S3_STORAGE_CLASS_GLACIER:  450,
			},
			priceList: MockProductPriceList,
			expectedOutput: map[string]float64{
				S3_STORAGE_CLASS_STANDARD: 0.25,
				S3_STORAGE_CLASS_GLACIER:  0.25,
			},
		},
		{
			name: "Test Second Tier",
			totalStorageClassSize: util.StorageClassSizeMap{
				S3_STORAGE_CLASS_STANDARD: 304,
				S3_STORAGE_CLASS_GLACIER:  2344234092384,
			},
			priceList: MockProductPriceList,
			expectedOutput: map[string]float64{
				S3_STORAGE_CLASS_STANDARD: 0.25,
				S3_STORAGE_CLASS_GLACIER:  0.20916070722464533,
			},
		},
	}
	for _, test := range tests {
		output := GetTierPriceList(test.totalStorageClassSize, test.priceList, float64(0))
		assert.Equal(t, test.expectedOutput, output)
	}
}
