package util

import (
	"sync"
	"time"
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

type BucketDTO struct {
	Name             string
	CreationDate     time.Time
	NbOfFiles        int64
	SizeOfBucket     float64
	LastUpdateDate   time.Time
	Cost             float64
	StorageClassSize StorageClassSizeMap
	Region           string
}

type StorageClassSize struct {
	SizeMap RegionsStorageMap
	Mutex   sync.Mutex
}
type RegionsStorageMap map[string]map[string]float64
type StorageClassSizeMap map[string]float64
