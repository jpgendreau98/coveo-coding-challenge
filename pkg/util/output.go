package util

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

func OutputData(buckets []*BucketDTO, options OutputOptions) error {
	output := make(map[string]interface{})
	if options.OrderByInc != "" {
		buckets = orderByInc(options.OrderByInc, buckets)

	} else if options.OrderByDec != "" {
		buckets = orderByDec(options.OrderByDec, buckets)

	}
	output["S3"] = buckets
	output["S3"] = applyOutputOptions(buckets, options)
	data, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return err
	}
	if options.FileOutput != "" {
		err := outputToFilePath(options.FileOutput, data)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(data))
	}
	return nil
}

func applyOutputOptions(data []*BucketDTO, outputOptions OutputOptions) map[string][]*BucketDTO {
	applyConversion := outputOptions.SizeConversion > 0
	applyGroupByRegion := outputOptions.GroupBy == "region"
	output := make(map[string][]*BucketDTO)
	for _, bucket := range data {
		if applyConversion {
			bucket = applySizeConversion(bucket, outputOptions.SizeConversion)
		}
		if applyGroupByRegion {
			output[bucket.Region] = append(output[bucket.Region], bucket)
		} else {
			output["Global"] = append(output["Global"], bucket)
		}
	}
	return output
}

func applySizeConversion(bucket *BucketDTO, sizeConversion float64) *BucketDTO {
	bucket.SizeOfBucket = bucket.SizeOfBucket / math.Pow(float64(1024), sizeConversion)
	for k, v := range bucket.StorageClassSize {
		bucket.StorageClassSize[k] = v / math.Pow(float64(1024), sizeConversion)
	}
	return bucket
}

func orderByInc(key string, data []*BucketDTO) []*BucketDTO {
	switch key {
	case "cost":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return cmp.Compare(a.Cost, b.Cost)
		})
	case "name":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return cmp.Compare(a.Name, b.Name)
		})
	case "size":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return cmp.Compare(a.SizeOfBucket, b.SizeOfBucket)
		})
	}
	return data
}

func orderByDec(key string, data []*BucketDTO) []*BucketDTO {
	switch key {
	case "cost":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return (cmp.Compare(a.Cost, b.Cost) * -1)
		})
	case "name":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return (cmp.Compare(a.Name, b.Name) * -1)
		})
	case "size":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return (cmp.Compare(a.SizeOfBucket, b.SizeOfBucket) * -1)
		})
	}
	return data
}

func outputToFilePath(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
