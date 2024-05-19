package main

import (
	"fmt"
	"projet-devops-coveo/pkg"
)

func main() {
	// p := tea.NewProgram(root.InitialBase())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
	fmt.Println("Fetching prices as of today...")
	svc := pkg.InitConnectionPricingList()
	var PriceList = make(pkg.MasterPriceList)
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
		PriceList[pkg.GetStorageClassNameBySky(storageClassSKU)] = decodedPriceList
	}

	fmt.Println("Price fetched Successfully!")
	fmt.Println("Starting the scrapping of S3 Buckets")

	options := pkg.CliOptions{
		// FilterByStorageClass: []string{"STANDARD"},
		ReturnEmptyBuckets: true,
		Regions:            []string{"ca-central-1", "us-east-1", "us-west-1"},
	}

	var allBuckets []*pkg.BucketDTO
	fs, err := pkg.InitConnection("ca-central-1")
	if err != nil {
		fmt.Println(err)
	}
	buckets := fs.ListDirectories(options)
	pkg.SortListBasedOnRegion(buckets)
	for _, region := range options.Regions {
		fs, err := pkg.InitConnection(region)
		if err != nil {
			fmt.Println(err)
		}
		regionBucket := pkg.GetBucketOfTheRegion(buckets, region)

		DTOBuckets := fs.GetObject(regionBucket, PriceList, options)
		fs.SetBucketPrices(DTOBuckets, PriceList)

		allBuckets = append(allBuckets, DTOBuckets...)
		buckets = pkg.RemoveScrappedBucketFromList(regionBucket, buckets)
	}
	fmt.Println("Buckets have been fetched successfuly!")
	pkg.OutputData(allBuckets, pkg.OutputOptions{
		// OrderByDec: "name",
	})

}
