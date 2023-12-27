package utils

import "math/rand"

func IsError() bool {
	if rand.Intn(100)%3 == 0 {
		return true
	}
	return false
}
