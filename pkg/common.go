package pkg

func TransformByteToGB(size int64) float64 {
	return float64(size) / (1024 * 1024 * 1024)
}
