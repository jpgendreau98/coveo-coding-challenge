package cmd

import (
	"fmt"
	"projet-devops-coveo/pkg"
)

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
