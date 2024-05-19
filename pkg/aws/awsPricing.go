package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"projet-devops-coveo/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/service/pricing"
)

type AwsPricing struct {
	Session *pricing.Pricing
}

type RegionSkuList map[string][]Product
type MasterPriceList map[string]ProductPriceList
type ProductPriceList map[string]PriceList

// Establish connections to aws pricing services.
func InitConnectionPricingList() *AwsPricing {
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	return &AwsPricing{
		Session: pricing.New(session),
	}
}

// Get skus of AmazonS3 storage products
func (ap *AwsPricing) GetSkusForRegions(regions []string) (RegionSkuList, error) {
	var regionSkuList = make(RegionSkuList)
	for _, region := range regions {
		results, err := ap.Session.ListPriceLists(&pricing.ListPriceListsInput{
			RegionCode:    aws.String(region),
			ServiceCode:   aws.String("AmazonS3"),
			CurrencyCode:  aws.String("USD"),
			EffectiveDate: aws.Time(time.Now()),
		})
		if err != nil {
			return nil, err
		}
		var priceslist []Product
		for _, priceList := range results.PriceLists {
			productPrices, err := ap.getProductWithArn(priceList.PriceListArn)
			if err != nil {
				return nil, err
			}
			priceslist = append(priceslist, productPrices...)
		}
		regionSkuList[region] = priceslist
	}
	return regionSkuList, nil

}

// Get PriceList of a product with an arn provided by AWS. Returns a list of products.
func (ap *AwsPricing) getProductWithArn(priceListArn *string) (list []Product, err error) {
	results, err := ap.Session.GetPriceListFileUrl(&pricing.GetPriceListFileUrlInput{
		FileFormat:   aws.String("json"),
		PriceListArn: priceListArn,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(aws.StringValue(results.Url))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	products := TEST{}
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}
	for _, product := range products.Products {
		if product.ProductFamily == "Storage" && product.Attributes.Operation == "" && product.Attributes.Usagetype != "TagStorage-TagHrs" {
			list = append(list, product)
		}
	}
	return list, nil
}

// Get Region Price price list with the sku list. Returns an Master price list of all prices in all regions.
func (ap *AwsPricing) GetRegionPriceList(regionSkuList RegionSkuList) MasterPriceList {
	regionMasterPriceList := make(MasterPriceList)
	for k, v := range regionSkuList {
		var productPriceList = make(ProductPriceList)
		for _, product := range v {
			productPrice, err := ap.GetPriceListWithSku(product.Sku)
			if err != nil {
				fmt.Println(err)
				continue
			}
			priceList, err := DecodePricingList(productPrice)
			if err != nil {
				fmt.Println(err)
				continue
			}
			productPriceList[product.Attributes.VolumeType] = priceList
		}
		regionMasterPriceList[k] = productPriceList
	}
	return regionMasterPriceList
}

// Get a price list with a sku.
func (ap *AwsPricing) GetPriceListWithSku(sku string) (*pricing.GetProductsOutput, error) {
	filters := []*pricing.Filter{{
		Field: aws.String("sku"),
		Value: aws.String(sku),
		Type:  aws.String("TERM_MATCH"),
	}}
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

// Decode a pricing list from AWS.
func DecodePricingList(productPrice *pricing.GetProductsOutput) (priceList PriceList, err error) {
	striout, _ := protocol.EncodeJSONValue(productPrice.PriceList[0], protocol.NoEscape)
	err = json.Unmarshal([]byte(striout), &priceList)
	if err != nil {
		return priceList, err
	}
	return priceList, nil
}

// Calculate the average price for a SKU. It will take the total size of all the storage class and return a map of average price per storage class.
func GetTierPriceList(totalStorageClassSize util.StorageClassSizeMap, priceList ProductPriceList, conversion float64) map[string]float64 {
	tierList := make(map[string]float64)
	for k, v := range totalStorageClassSize {
		priceListForSku := priceList[GetStorageClassType(k)]
		price, err := getPriceForSize(TransformSizeToGB(v, conversion), priceListForSku)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tierList[k] = price
	}
	return tierList

}

// Function to calculate the average price per storage class based on total size of storage class.
func getPriceForSize(sizeGB float64, priceListForSku PriceList) (price float64, err error) {
	var totalPrice float64
	var tempSize = sizeGB
	for _, j := range priceListForSku.Terms.OnDemand {
		for _, l := range j.PriceDimensions {
			bRange, err := strconv.ParseFloat(l.BeginRange, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			unitPrice, err := strconv.ParseFloat(l.PricePerUnit.Usd, 32)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if l.EndRange != "Inf" {
				eRange, err := strconv.ParseFloat(l.EndRange, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if (eRange - bRange) >= tempSize {
					totalPrice += tempSize * unitPrice
					break
				} else {
					totalPrice += (eRange - bRange) * unitPrice
					tempSize -= (eRange - bRange)
				}
			} else {
				totalPrice += tempSize * unitPrice
			}
		}
	}
	return totalPrice / sizeGB, nil
}

// Help function to convert between AWS Bucket Storage class and AWS Price liste Storage Class
func GetStorageClassType(volumeType string) string {
	switch volumeType {
	case S3_STORAGE_CLASS_GLACIER_IR:
		return "Glacier Instant Retrieval"
	case S3_STORAGE_CLASS_STANDARD:
		return "Standard"
	case S3_STORAGE_CLASS_INTELLIGENT_TIERING:
		return "Intelligent-Tiering Frequent Access"
	case S3_STORAGE_CLASS_GLACIER:
		return "Amazon Glacier"
	case S3_STORAGE_CLASS_REDUCED_REDUNDANCY:
		return "Reduced Redundancy"
	case S3_STORAGE_CLASS_STANDARD_IA:
		return "Standard - Infrequent Access"
	default:
		return S3_STORAGE_CLASS_STANDARD
	}
}
