package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		randCharsetIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(fmt.Sprintf("Не удалось сгенерировать случайное число: %s", err.Error()))
		}
		b[i] = charset[randCharsetIndex.Int64()]
	}

	return string(b)
}

func RandomInt(min, max int) int {
	diff := big.NewInt(int64(max - min + 1))
	randInt, err := rand.Int(rand.Reader, diff)
	if err != nil {
		panic(fmt.Sprintf("Не удалось сгенерировать случайное число: %s", err.Error()))
	}
	return int(randInt.Int64()) + min
}
