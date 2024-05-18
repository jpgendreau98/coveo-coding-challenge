package main

import (
	"fmt"
	"projet-devops-coveo/pkg"
)

type Price struct {
	Terms Terms `json:"terms"`
}

type Terms struct {
	OnDemand map[string]map[string]interface{} `json:"onDemand"`
}

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
	fs, err := pkg.InitConnection()
	if err != nil {
		fmt.Println(err)
	}

	buckets := fs.ListDirectories()
	DTOBuckets := fs.GetObjectMetadata(buckets, PriceList)
	fs.GetBucketPrices(DTOBuckets, PriceList)
	for _, bucket := range DTOBuckets {
		fmt.Printf("%+v \n", bucket)
	}
	fmt.Println("Buckets have been fetched successfuly!")

}
