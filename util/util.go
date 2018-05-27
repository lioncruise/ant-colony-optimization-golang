package util

import "math"

// 在Golang 1.10中提供了math.Round方法，不需要自己实现
func Round(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}