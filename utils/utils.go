package utils

import (
	"math"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	newVal := math.Round(val*ratio) / ratio
	return newVal
}
