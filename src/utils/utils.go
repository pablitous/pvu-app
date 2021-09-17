package utils

import "math/rand"

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
