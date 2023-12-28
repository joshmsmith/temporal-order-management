package utils

import "math/rand"

func IsError() bool {
	// throw errors 1 time out of 100
	if rand.Intn(100) > 99 {
		return true
	}
	return false
}
