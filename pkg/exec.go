package pkg

import (
	"fmt"
	"sync"

	"projet-devops-coveo/pkg/aws"
	"projet-devops-coveo/pkg/util"

	"go.uber.org/ratelimit"
)

func RunS3Command(options *util.CliOptions) error {
	//Start the ratelimiter
	limiter := ratelimit.New(options.RateLimit)
	//Fetching the price of the day.
	fmt.Println("Fetching prices as of today...")
	priceList, err := fetchPrices(*options)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Price fetched Successfully!")

	fmt.Println("Starting the scrapping of S3 Buckets")

	// Init GlobalStorageMap to use it on all s3 regions
	globalStorageClassSize := initRegionStorageMap(options.Regions)
	var allBuckets []*util.BucketDTO
	// Create initial connection for scrapping of all the buckets since this call is regionless
	fs, err := aws.InitConnection(options.Regions[0], *options, globalStorageClassSize, limiter)
	if err != nil {
		fmt.Println(err)
	}
	// Filter buckets with the filter given by user
	buckets := fs.GetBucketsFiltered()
	// We have to sort the liste of buckets for increase performance for search functions
	aws.SortListBasedOnRegion(buckets)
	//Since all the region
	bucketChan := make(chan []*util.BucketDTO, len(buckets))
	wg := new(sync.WaitGroup)
	for _, region := range options.Regions {
		wg.Add(1)
		//Init a new connection with the region
		fs, err := aws.InitConnection(region, *options, globalStorageClassSize, limiter)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//Return all the bucket in the region
		regionBucket := aws.GetBucketOfTheRegion(buckets, region)
		//Using the sorted lists from earlier, the search is way faster to find the index of the buckets
		buckets = aws.RemoveScrappedBucketFromList(regionBucket, buckets)
		//Starting multi-threading on the scrap of objects.
		go fs.ListObjectsInBucket(regionBucket, region, priceList, wg, bucketChan)

	}
	wg.Wait()
	close(bucketChan)
	for bucket := range bucketChan {
		allBuckets = append(allBuckets, bucket...)
	}
	// Set Bucket cost with all the information gathered.
	fs.SetBucketCost(allBuckets, priceList)
	fmt.Println("Buckets have been fetched successfuly!")
	fmt.Println("Printing data...")
	util.OutputData(allBuckets, *options.OutputOptions)
	return nil
}

func initRegionStorageMap(regions []string) *util.StorageClassSize {
	var globalStorageClassSize = &util.StorageClassSize{
		SizeMap: make(util.RegionsStorageMap),
	}
	for _, region := range regions {
		globalStorageClassSize.SizeMap[region] = make(map[string]float64)
	}
	return globalStorageClassSize
}

func fetchPrices(options util.CliOptions) (aws.MasterPriceList, error) {
	svc := aws.InitConnectionPricingList()
	regionSkuList, err := svc.GetSkusForRegions(options.Regions)
	if err != nil {
		return nil, err
	}
	masterPriceList := svc.GetRegionPriceList(regionSkuList)
	return masterPriceList, nil
}
