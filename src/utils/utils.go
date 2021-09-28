package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func RandFloats(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

func AddRandomSleep(si float64, se float64) bool {
	rand.Seed(time.Now().UnixNano())
	n := RandFloats(si, se)
	fmt.Printf("Waiting %f seconds...\n", n)
	time.Sleep(time.Duration(n) * time.Second)
	return true
}
