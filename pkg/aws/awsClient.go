package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/ratelimit"
)

type AwsInterface interface{}

type AwsClient struct {
	s3      *s3.Client
	limiter ratelimit.Limiter
	pricing *pricing.Client
	ctx     context.Context
}

func NewAwsClient(region string, limiter ratelimit.Limiter) (*AwsClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	return &AwsClient{
		s3: s3.New(s3.Options{
			Region:      region,
			Credentials: cfg.Credentials,
		}),
		ctx:     context.Background(),
		limiter: limiter,
		pricing: pricing.New(pricing.Options{
			Region:      "us-east-1",
			Credentials: cfg.Credentials,
		}),
	}, nil
}

func (a *AwsClient) ListBuckets(input *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	a.limiter.Take()
	return a.s3.ListBuckets(a.ctx, input, optFns...)
}

func (a *AwsClient) GetBucketLocation(params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (string, error) {
	a.limiter.Take()
	output, err := a.s3.GetBucketLocation(a.ctx, params, optFns...)
	if err != nil {
		return "", err
	}
	return GetBucketLocationConstant(output.LocationConstraint), nil
}

func (a *AwsClient) ListObjectsV2(params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	a.limiter.Take()
	return a.s3.ListObjectsV2(a.ctx, params, optFns...)
}

func (a *AwsClient) ListPriceLists(params *pricing.ListPriceListsInput, optFns ...func(*pricing.Options)) (*pricing.ListPriceListsOutput, error) {
	a.limiter.Take()
	return a.pricing.ListPriceLists(a.ctx, params, optFns...)
}

func (a *AwsClient) GetPriceListFileUrl(params *pricing.GetPriceListFileUrlInput, optFns ...func(*pricing.Options)) (*pricing.GetPriceListFileUrlOutput, error) {
	a.limiter.Take()
	return a.pricing.GetPriceListFileUrl(a.ctx, params, optFns...)
}

func (a *AwsClient) GetProducts(params *pricing.GetProductsInput, optFns ...func(*pricing.Options)) (*pricing.GetProductsOutput, error) {
	a.limiter.Take()
	return a.pricing.GetProducts(a.ctx, params, optFns...)
}
