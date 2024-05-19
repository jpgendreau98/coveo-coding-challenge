package util

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

func OutputData(buckets []*BucketDTO, options OutputOptions) error {
	output := make(map[string]interface{})

	if options.OrderByInc != "" {
		buckets = orderByInc(options.OrderByInc, buckets)

	} else if options.OrderByDec != "" {
		buckets = orderByDec(options.OrderByInc, buckets)

	}
	output["S3"] = buckets
	if options.GroupBy != "" {
		switch options.GroupBy {
		case "region":
			output["S3"] = groupByRegion(buckets)
		}
	}
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

func groupByRegion(data []*BucketDTO) map[string][]*BucketDTO {
	groupByOutput := make(map[string][]*BucketDTO)
	for _, bucket := range data {
		groupByOutput[bucket.Region] = append(groupByOutput[bucket.Region], bucket)
	}
	return groupByOutput
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
			return cmp.Compare(a.Cost, b.Cost)
		})
		slices.Reverse(data)
	case "name":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return cmp.Compare(a.Name, b.Name)
		})
		slices.Reverse(data)
	case "size":
		slices.SortStableFunc(data, func(a, b *BucketDTO) int {
			return cmp.Compare(a.SizeOfBucket, b.SizeOfBucket)
		})
		slices.Reverse(data)
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
