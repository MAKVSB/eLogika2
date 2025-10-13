package utils

import "math"

func RoundToEven(f float64, decimals int) float64 {
	if decimals < 0 {
		return f
	}
	factor := math.Pow(10, float64(decimals))
	return math.RoundToEven(f*factor) / factor
}

func ToPercentage(f uint) float64 {
	return float64(f) / float64(100)
}
