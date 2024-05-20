package util

import (
	"math"
	"sync"
)

type CliOptions struct {
	FilterByName         []string
	FilterByStorageClass []string
	OmitEmpty            bool
	Regions              []string
	OutputOptions        *OutputOptions
	RateLimit            int
	Threading            int
}

type OutputOptions struct {
	GroupBy        string
	OrderByDec     string
	OrderByInc     string
	FileOutput     string
	SizeConversion float64
}

type StorageClassSize struct {
	SizeMap RegionsStorageMap
	Mutex   sync.Mutex
}
type RegionsStorageMap map[string]map[string]float64
type StorageClassSizeMap map[string]float64

func (smap *RegionsStorageMap) ApplyConversion(conversion float64) {
	for _, v := range *smap {
		for k, s := range v {
			v[k] = s / math.Pow(float64(1024), conversion)
		}
	}
}
