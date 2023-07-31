package zkp_auth

import "math"

func Pow(a, b int64) int64 {
	return int64(math.Pow(float64(a), float64(b)))
}

func Mod(a, m int64) int64 {
	r := a % m
	if r < 0 {
		return r + m
	}
	return r
}

func calculateCommitment(exp int64) (int64, int64) {
	return Mod(Pow(g, exp), p), Mod(Pow(h, exp), p)
}
