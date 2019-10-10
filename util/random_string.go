package util

import (
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789~!@#$%^&*()_+")
var chars = []rune("abcdefghijklmnopqrstuvwxyz")
var digitsLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

var digits = []rune("1234567890")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func RandomDigitAndLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(digitsLetters))]
	}
	return string(b)
}

func RandomWord(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(chars))]
	}
	return string(b)
}

func RandEmail() string {
	return RandomWord(6) + "@" + RandomWord(4) + "." + RandomWord(3)
}
