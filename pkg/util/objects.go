package util

import (
	"sync"
	"time"
)

type BucketDTO struct {
	Name             string
	CreationDate     time.Time
	NbOfFiles        int64
	SizeOfBucket     float64
	LastUpdateDate   time.Time
	Cost             float64
	StorageClassSize StorageClassSize
	Region           string
}

type StorageClassSize struct {
	SizeMap StorageClassSizeMap
	Mutex   sync.Mutex
}

type StorageClassSizeMap map[string]float64
