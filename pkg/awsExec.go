package pkg

import (
	"sync"
	"time"

	"projet-devops-coveo/pkg/aws"
	"projet-devops-coveo/pkg/util"

	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

func RunS3Command(options *util.CliOptions) error {
	//Start the ratelimiter
	limiter := ratelimit.New(options.RateLimit)
	awsClient, err := aws.NewAwsClient(options.Regions[0], limiter)
	//Fetching the price of the day.
	logrus.Info("Fetching prices as of today...")
	priceList, err := fetchPrices(awsClient, *options)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Price fetched Successfully!")

	logrus.Info("Starting the scrapping of S3 Buckets")
	start := time.Now()
	// Init GlobalStorageMap to use it on all s3 regions
	globalStorageClassSize := initRegionStorageMap(options.Regions)
	var allBuckets []util.CloudFilesystem
	// Create initial connection for scrapping of all the buckets since this call is regionless
	fs, err := aws.InitConnection(options.Regions[0], *options, globalStorageClassSize, limiter)
	if err != nil {
		logrus.Error(err)
	}
	// Filter buckets with the filter given by user (Filter by name and Filter by region)
	buckets := fs.GetBucketsFiltered()
	// We have to sort the list of buckets for increase performance for search functions
	aws.SortListBasedOnRegion(buckets)
	//Since the sdk of Go doesn't let you scrap a bucket which is not in the region of the config,
	//we have to loop on all the wanted regions to be able to scrap all the buckets.
	bucketChan := make(chan []util.CloudFilesystem, len(buckets))
	wg := new(sync.WaitGroup)
	for _, region := range options.Regions {
		wg.Add(1)
		//Init a new connection with the region
		fs, err := aws.InitConnection(region, *options, globalStorageClassSize, limiter)
		if err != nil {
			logrus.Error(err)
			continue
		}
		//Return all the bucket in the region
		regionBucket := aws.GetBucketsOfRegion(buckets, region)
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
	logrus.Info("Buckets have been fetched successfuly!")
	logrus.Info("Execution Time: ", time.Since(start))
	logrus.Info("Printing data...")
	//Print Data
	util.OutputData(allBuckets, *options.OutputOptions, globalStorageClassSize.SizeMap)
	logrus.Info("Done!")
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

func fetchPrices(awsClient aws.AwsInterface, options util.CliOptions) (aws.MasterPriceList, error) {
	//Init connection to AWS pricing services
	svc := aws.InitConnectionPricingList(awsClient)
	//Get a list with all the skus for Amazon S3 product grouped by region
	regionSkuList, err := svc.GetSkusForRegions(options.Regions)
	if err != nil {
		return nil, err
	}
	//Create a price list with all the different prices for the wanted regions
	masterPriceList := svc.GetRegionPriceList(regionSkuList)
	return masterPriceList, nil
}
