package util

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

func OutputData(buckets []CloudFilesystem, options OutputOptions, gloablStorageClass RegionsStorageMap) error {
	output := make(map[string]interface{})
	if options.OrderByInc != "" {
		buckets = orderByInc(options.OrderByInc, buckets)

	} else if options.OrderByDec != "" {
		buckets = orderByDec(options.OrderByDec, buckets)

	}
	output["S3"] = buckets
	output["S3"] = applyOutputOptions(buckets, options)
	gloablStorageClass.ApplyConversion(options.SizeConversion)
	output["S3-Stats"] = gloablStorageClass
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

func applyOutputOptions(data []CloudFilesystem, outputOptions OutputOptions) map[string][]CloudFilesystem {
	applyConversion := outputOptions.SizeConversion > 0
	applyGroupByRegion := outputOptions.GroupBy == "region"
	output := make(map[string][]CloudFilesystem)
	if applyConversion || applyGroupByRegion {
		for _, bucket := range data {
			if applyConversion {
				bucket.ApplySizeConversion(outputOptions.SizeConversion)
			}
			if applyGroupByRegion {
				output[bucket.GetRegion()] = append(output[bucket.GetRegion()], bucket)
			} else {
				output["Global"] = append(output["Global"], bucket)
			}
		}
	}
	return output
}

func orderByInc(key string, data []CloudFilesystem) []CloudFilesystem {
	switch key {
	case "cost":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return cmp.Compare(a.GetCost(), b.GetCost())
		})
	case "name":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return cmp.Compare(a.GetName(), b.GetName())
		})
	case "size":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return cmp.Compare(a.GetSizeOfBucket(), b.GetSizeOfBucket())
		})
	}
	return data
}

func orderByDec(key string, data []CloudFilesystem) []CloudFilesystem {
	switch key {
	case "cost":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return (cmp.Compare(a.GetCost(), b.GetCost()) * -1)
		})
	case "name":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return (cmp.Compare(a.GetName(), b.GetName()) * -1)
		})
	case "size":
		slices.SortStableFunc(data, func(a, b CloudFilesystem) int {
			return (cmp.Compare(a.GetSizeOfBucket(), b.GetSizeOfBucket()) * -1)
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
