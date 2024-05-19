package cmd

import (
	"fmt"
	"projet-devops-coveo/pkg"
	"projet-devops-coveo/pkg/util"

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

	RETURNS_EMTPY             = "omit-empty"
	RETURNS_EMTPY_DESCRIPTION = "Omit empty buckets"
	RETURNS_EMTPY_DEFAULT     = false

	RATE_LIMIT             = "ratelimit"
	RATE_LIMIT_DESCRIPTION = "Choose the rate limit on the S3 services (Max Get = 5500 per second)"
	RATE_LIMIT_DEFAULT     = 5400

	THREADING             = "threading"
	THREADING_DESCRIPTION = "Choose the number of concurrent task for performance, This is multiply per number of regions that you want to scrap"
	THREADING_DEFAULT     = 400

	SIZE_CONV             = "display-size"
	SIZE_CONV_DESCRIPTION = "Display the size in: [by, kb, mb, gb, tb]"
	SIZE_CONV_DEFAULT     = "by"

	FILTER_BY_NAME             = "name"
	FILTER_BY_NAME_DESCRIPTION = "Select multiples bucket name to."

	FILTER_BY_STORAGE_CLASS             = "storage-class"
	FILTER_BY_STORAGE_CLASS_DESCRIPTION = "Select multiples storage class to filter bucket contents. Supported: [STANDARD, REDUCED_REDUNDANCY, GLACIER, STANDARD_IA, INTELLIGENT_TIERING, DEEP_ARCHIVE, GLACIER_IR]"

	BUCKET_REGIONS             = "regions"
	BUCKET_REGIONS_DESCRIPTION = "Regions in which the bucket are created"

	OUTPUT             = "output"
	OUTPUT_DESCRIPTION = "Output to a file (Enter the file name)"
)

func NewS3Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "aws-s3",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := &util.CliOptions{
				Regions:              viper.GetStringSlice(BUCKET_REGIONS),
				FilterByName:         viper.GetStringSlice(FILTER_BY_NAME),
				OmitEmpty:            viper.GetBool(RETURNS_EMTPY),
				FilterByStorageClass: viper.GetStringSlice(FILTER_BY_STORAGE_CLASS),
				RateLimit:            viper.GetInt(RATE_LIMIT),
				Threading:            viper.GetInt(THREADING),
				OutputOptions: &util.OutputOptions{
					GroupBy:        viper.GetString(GROUP_BY),
					OrderByInc:     viper.GetString(ORDER_BY_INC),
					OrderByDec:     viper.GetString(ORDER_BY_DEC),
					FileOutput:     viper.GetString(OUTPUT),
					SizeConversion: float64(getSizeConstant(viper.GetString(SIZE_CONV))),
				},
			}
			err := pkg.RunS3Command(options)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(ORDER_BY_INC, ORDER_BY_INC_DEFAULT, ORDER_BY_INC_DESCRIPTION)
	cmd.Flags().String(ORDER_BY_DEC, ORDER_BY_DEC_DEFAULT, ORDER_BY_DEC_DESCRIPTION)
	cmd.Flags().String(GROUP_BY, GROUP_BY_DEFAULT, GROUP_BY_DESCRIPTION)
	cmd.Flags().String(SIZE_CONV, SIZE_CONV_DEFAULT, SIZE_CONV_DESCRIPTION)
	cmd.Flags().Int(RATE_LIMIT, RATE_LIMIT_DEFAULT, RATE_LIMIT_DESCRIPTION)
	cmd.Flags().Int(THREADING, THREADING_DEFAULT, THREADING_DESCRIPTION)
	cmd.Flags().String(OUTPUT, "", OUTPUT_DESCRIPTION)
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

func getSizeConstant(size string) int {
	switch size {
	case "by":
		return util.SIZE_CONV_BY
	case "kb":
		return util.SIZE_CONV_KB
	case "mb":
		return util.SIZE_CONV_MB
	case "gb":
		return util.SIZE_CONV_GB
	case "tb":
		return util.SIZE_CONV_TB
	}
	return util.SIZE_CONV_BY
}
