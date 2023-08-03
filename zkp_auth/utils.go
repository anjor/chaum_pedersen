package zkp_auth

import (
	"math"
	"math/rand"
)

func pow(a, b int64) int64 {
	return int64(math.Pow(float64(a), float64(b)))
}

func mod(a, m int64) int64 {
	r := a % m
	if r < 0 {
		return r + m
	}
	return r
}

func calculateCommitment(exp int64) (int64, int64) {
	return mod(pow(g, exp), p), mod(pow(h, exp), p)
}

func generateRandom() int64 {
	return int64(rand.Intn(10) + 1)
}
