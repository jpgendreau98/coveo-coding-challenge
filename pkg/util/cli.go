package util

type CliOptions struct {
	FilterByName         []string
	FilterByStorageClass []string
	ReturnEmptyBuckets   bool
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
