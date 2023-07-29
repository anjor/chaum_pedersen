package zkp_auth

import "math"

func Pow(a, b int64) int64 {
	return int64(math.Pow(float64(a), float64(b)))
}

func Mod(a, m int64) int64 {
	r := int64(math.Mod(float64(a), float64(m)))
	if r < 0 {
		return r + m
	}
	return r
}
