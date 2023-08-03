package zkp_auth

import (
	"testing"
)

func TestPow(t *testing.T) {
	testCases := []struct {
		a, b, expected int64
	}{
		{2, 3, 8},   // 2^3 = 8
		{5, 0, 1},   // 5^0 = 1
		{0, 5, 0},   // 0^5 = 0
		{-3, 2, 9},  // (-3)^2 = 9
		{-2, 3, -8}, // (-2)^3 = -8
		{10, -1, 0}, // 10^(-1) = 0.1 (rounded down to 0 due to int64 return type)
	}

	for _, tc := range testCases {
		result := pow(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("pow(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestMod(t *testing.T) {
	testCases := []struct {
		a, m, expected int64
	}{
		{5, 3, 2},   // 5 % 3 = 2
		{-5, 3, 1},  // (-5) % 3 = 1
		{10, 5, 0},  // 10 % 5 = 0
		{17, 7, 3},  // 17 % 7 = 3
		{0, 8, 0},   // 0 % 8 = 0
		{-10, 6, 2}, // (-10) % 6 = 2
	}

	for _, tc := range testCases {
		result := mod(tc.a, tc.m)
		if result != tc.expected {
			t.Errorf("mod(%d, %d) = %d; expected %d", tc.a, tc.m, result, tc.expected)
		}
	}
}

func TestCalculateCommitment(t *testing.T) {
	testCases := []struct {
		exp        int64
		expectedY1 int64
		expectedY2 int64
	}{
		{0, 1, 1},
		{1, 4, 9},
		{2, 16, 12},
		{3, 18, 16},
	}

	for _, tc := range testCases {
		y1, y2 := calculateCommitment(tc.exp)
		if y1 != tc.expectedY1 || y2 != tc.expectedY2 {
			t.Errorf("calculateCommitment(%d) returned y1: %d, y2: %d; expected y1: %d, y2: %d",
				tc.exp, y1, y2, tc.expectedY1, tc.expectedY2)
		}
	}
}
