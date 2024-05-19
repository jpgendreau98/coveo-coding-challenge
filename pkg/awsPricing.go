package pkg

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/service/pricing"
)

type AwsPricing struct {
	Session *pricing.Pricing
}

type MasterPriceList map[string]PriceList

func InitConnectionPricingList() *AwsPricing {
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	return &AwsPricing{
		Session: pricing.New(session),
	}
}

func (ap *AwsPricing) GetPriceListWithSku(sku string) (*pricing.GetProductsOutput, error) {
	filters := []*pricing.Filter{{
		Field: aws.String("sku"),
		Value: aws.String(sku),
		Type:  aws.String("TERM_MATCH"),
	},
	}
	input := pricing.GetProductsInput{
		Filters:     filters,
		ServiceCode: aws.String("AmazonS3"),
	}
	productPrice, err := ap.Session.GetProducts(&input)
	if err != nil {
		return nil, err
	}
	return productPrice, nil
}

func DecodePricingList(productPrice *pricing.GetProductsOutput) (priceList PriceList, err error) {
	striout, _ := protocol.EncodeJSONValue(productPrice.PriceList[0], protocol.NoEscape)
	err = json.Unmarshal([]byte(striout), &priceList)
	if err != nil {
		return priceList, err
	}
	return priceList, nil
}

func GetTierPriceList(totalStorageClassSize StorageClassSize, priceList MasterPriceList) map[string]float64 {
	tierList := make(map[string]float64)
	for k, v := range totalStorageClassSize {
		priceListForSku := priceList[k]
		//Transform size in GB
		sizeGB := float64(v) / (1024 * 1024 * 1024)
		price, err := getPriceForSize(sizeGB, priceListForSku)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tierList[k] = price
	}
	return tierList

}

func getPriceForSize(sizeGB float64, priceListForSku PriceList) (price float64, err error) {
	for _, j := range priceListForSku.Terms.OnDemand {
		for _, l := range j.PriceDimensions {
			bRange, err := strconv.ParseFloat(l.BeginRange, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var eRange float64
			if l.EndRange != "Inf" {
				eRange, err = strconv.ParseFloat(l.EndRange, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if sizeGB > bRange && sizeGB < eRange {
					unitPrice, err := strconv.ParseFloat(l.PricePerUnit.Usd, 32)
					if err != nil {
						fmt.Println(err)
						continue
					}
					return unitPrice, nil
				}
			} else {
				if sizeGB > bRange {
					unitPrice, err := strconv.ParseFloat(l.PricePerUnit.Usd, 32)
					if err != nil {
						fmt.Println(err)
						continue
					}
					return unitPrice, nil
				}
			}

		}
	}
	return 0, fmt.Errorf("Error while fetching tier price for the size of the storage class")
}

func GetStorageClassSku(storageClass string) string {
	switch storageClass {
	case "STANDARD":
		return S3_SKU_VOL_STANDARD
	case "REDUCED_REDUNDANCY":
		return S3_SKU_VOL_REDUCED_REDUNDANCY
	case "GLACIER":
		return S3_SKU_VOL_AMAZON_GLACIER
	case "STANDARD_IA":
		return S3_SKU_VOL_STANDARD_INFREQUENT_ACCESS
	case "ONEZONE_IA":
		return S3_SKU_VOL_STANDARD
	case "INTELLIGENT_TIERING":
		return S3_SKU_VOL_INTELLIGENT_TIERING_FREQUENT_ACCESS
	case "DEEP_ARCHIVE":
		return S3_SKU_VOL_GLACIER_DEEP_ARCHIVE
	case "OUTPOSTS":
		return S3_SKU_VOL_STANDARD
	case "GLACIER_IR":
		return S3_SKU_VOL_GLACIER_INSTANT_RETRIEVAL
	case "SNOW":
		return S3_SKU_VOL_STANDARD
	case "EXPRESS_ONEZONE":
		return S3_SKU_VOL_STANDARD
	default:
		return S3_SKU_VOL_STANDARD
	}
}

func GetStorageClassNameBySky(sku string) string {
	switch sku {
	case S3_SKU_VOL_STANDARD:
		return "STANDARD"
	case S3_SKU_VOL_REDUCED_REDUNDANCY:
		return "REDUCED_REDUNDANCY"
	case S3_SKU_VOL_AMAZON_GLACIER:
		return "GLACIER"
	case S3_SKU_VOL_STANDARD_INFREQUENT_ACCESS:
		return "STANDARD_IA"
	case S3_SKU_VOL_INTELLIGENT_TIERING_FREQUENT_ACCESS:
		return "INTELLIGENT_TIERING"
	case S3_SKU_VOL_GLACIER_DEEP_ARCHIVE:
		return "DEEP_ARCHIVE"
	case S3_SKU_VOL_GLACIER_INSTANT_RETRIEVAL:
		return "GLACIER_IR"
	default:
		return "STANDARD"
	}
}
