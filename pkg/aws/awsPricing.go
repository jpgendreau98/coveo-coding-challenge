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

func InitConnectionPricingList() *AwsPricing {
	session, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	return &AwsPricing{
		Session: pricing.New(session),
	}
}

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
		pricelist, err := ap.getPriceListWithArn(results.PriceLists)
		if err != nil {
			return nil, err
		}
		regionSkuList[region] = pricelist
	}
	return regionSkuList, nil

}

func (ap *AwsPricing) getPriceListWithArn(pricelists []*pricing.PriceList) (list []Product, err error) {
	for _, priceList := range pricelists {
		results, err := ap.Session.GetPriceListFileUrl(&pricing.GetPriceListFileUrlInput{
			FileFormat:   aws.String("json"),
			PriceListArn: priceList.PriceListArn,
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
	}
	return list, nil
}

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

func DecodePricingList(productPrice *pricing.GetProductsOutput) (priceList PriceList, err error) {
	striout, _ := protocol.EncodeJSONValue(productPrice.PriceList[0], protocol.NoEscape)
	err = json.Unmarshal([]byte(striout), &priceList)
	if err != nil {
		return priceList, err
	}
	return priceList, nil
}

func GetTierPriceList(totalStorageClassSize util.StorageClassSizeMap, priceList ProductPriceList) map[string]float64 {
	tierList := make(map[string]float64)
	for k, v := range totalStorageClassSize {
		priceListForSku := priceList[GetStorageClassType(k)]
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
				if (eRange - bRange) > tempSize {
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
		return S3_STORAGE_CLASS_STANDARD
	case S3_SKU_VOL_REDUCED_REDUNDANCY:
		return S3_STORAGE_CLASS_REDUCED_REDUNDANCY
	case S3_SKU_VOL_AMAZON_GLACIER:
		return S3_STORAGE_CLASS_GLACIER
	case S3_SKU_VOL_STANDARD_INFREQUENT_ACCESS:
		return S3_STORAGE_CLASS_STANDARD_IA
	case S3_SKU_VOL_INTELLIGENT_TIERING_FREQUENT_ACCESS:
		return S3_STORAGE_CLASS_INTELLIGENT_TIERING
	case S3_SKU_VOL_GLACIER_DEEP_ARCHIVE:
		return S3_STORAGE_CLASS_DEEP_ARCHIVE
	case S3_SKU_VOL_GLACIER_INSTANT_RETRIEVAL:
		return S3_STORAGE_CLASS_GLACIER_IR
	default:
		return S3_STORAGE_CLASS_STANDARD
	}
}
