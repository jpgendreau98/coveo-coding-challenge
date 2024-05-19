package cmd

import (
	"fmt"
	"projet-devops-coveo/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ORDER_BY_INC             = "order-by-inc"
	ORDER_BY_INC_DESCRIPTION = "Order by : name , price, size, storage-class"
	ORDER_BY_INC_DEFAULT     = ""

	ORDER_BY_DEC             = "order-by-dec"
	ORDER_BY_DEC_DESCRIPTION = "Order by : name , price, size, storage-class"
	ORDER_BY_DEC_DEFAULT     = ""

	GROUP_BY             = "group-by"
	GROUP_BY_DESCRIPTION = "Supported: [region]"
	GROUP_BY_DEFAULT     = ""

	RETURNS_EMTPY             = "return-emtpy-buckets"
	RETURNS_EMTPY_DESCRIPTION = "Omit empty buckets"
	RETURNS_EMTPY_DEFAULT     = true

	FILTER_BY_NAME             = "name"
	FILTER_BY_NAME_DESCRIPTION = "Select multiples bucket name to."

	FILTER_BY_STORAGE_CLASS             = "storage-class"
	FILTER_BY_STORAGE_CLASS_DESCRIPTION = "Select multiples storage class to filter bucket contents. Supported: [STANDARD, REDUCED_REDUNDANCY, GLACIER, STANDARD_IA, INTELLIGENT_TIERING, DEEP_ARCHIVE, GLACIER_IR]"

	BUCKET_REGIONS             = "regions"
	BUCKET_REGIONS_DESCRIPTION = "Regions in which the bucket are created"
)

func NewS3Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "aws-s3",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := &pkg.CliOptions{
				Regions:              viper.GetStringSlice(BUCKET_REGIONS),
				GroupBy:              viper.GetString(GROUP_BY),
				OrderByINC:           viper.GetString(ORDER_BY_INC),
				OrderByDEC:           viper.GetString(ORDER_BY_DEC),
				FilterByName:         viper.GetStringSlice(FILTER_BY_NAME),
				ReturnEmptyBuckets:   viper.GetBool(RETURNS_EMTPY),
				FilterByStorageClass: viper.GetStringSlice(FILTER_BY_STORAGE_CLASS),
			}
			err := RunS3Command(options, &pkg.OutputOptions{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(ORDER_BY_INC, ORDER_BY_INC_DEFAULT, ORDER_BY_INC_DESCRIPTION)
	cmd.Flags().String(ORDER_BY_DEC, ORDER_BY_DEC_DEFAULT, ORDER_BY_DEC_DESCRIPTION)
	cmd.Flags().String(GROUP_BY, GROUP_BY_DEFAULT, GROUP_BY_DESCRIPTION)
	cmd.Flags().Bool(RETURNS_EMTPY, RETURNS_EMTPY_DEFAULT, RETURNS_EMTPY_DESCRIPTION)
	cmd.Flags().StringSlice(FILTER_BY_NAME, nil, FILTER_BY_NAME_DESCRIPTION)
	cmd.Flags().StringSlice(BUCKET_REGIONS, []string{"ca-central-1"}, BUCKET_REGIONS_DESCRIPTION)
	cmd.Flags().StringSlice(FILTER_BY_STORAGE_CLASS, nil, FILTER_BY_STORAGE_CLASS_DESCRIPTION)
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return cmd
}

func RunS3Command(options *pkg.CliOptions, outputOptions *pkg.OutputOptions) error {
	fmt.Println("Fetching prices as of today...")
	priceList := fetchPrices()
	fmt.Println("Price fetched Successfully!")

	fmt.Println("Starting the scrapping of S3 Buckets")

	var allBuckets []*pkg.BucketDTO
	fs, err := pkg.InitConnection("ca-central-1")
	if err != nil {
		fmt.Println(err)
	}
	buckets := fs.ListDirectories(*options)
	pkg.SortListBasedOnRegion(buckets)
	for _, region := range options.Regions {
		fs, err := pkg.InitConnection(region)
		if err != nil {
			fmt.Println(err)
		}
		regionBucket := pkg.GetBucketOfTheRegion(buckets, region)

		DTOBuckets := fs.GetObject(regionBucket, priceList, *options)
		fs.SetBucketPrices(DTOBuckets, priceList)

		allBuckets = append(allBuckets, DTOBuckets...)
		buckets = pkg.RemoveScrappedBucketFromList(regionBucket, buckets)
	}
	fmt.Println("Buckets have been fetched successfuly!")
	pkg.OutputData(allBuckets, pkg.OutputOptions{
		// OrderByDec: "name",
	})
	return nil
}

func fetchPrices() pkg.MasterPriceList {
	svc := pkg.InitConnectionPricingList()
	var priceList = make(pkg.MasterPriceList)
	for _, storageClassSKU := range pkg.StorageClassesSKU {
		result, err := svc.GetPriceListWithSku(storageClassSKU)
		if err != nil {
			fmt.Println(err)
			continue
		}
		decodedPriceList, err := pkg.DecodePricingList(result)
		if err != nil {
			fmt.Println(err)
			continue
		}
		priceList[pkg.GetStorageClassNameBySky(storageClassSKU)] = decodedPriceList
	}
	return priceList
}
