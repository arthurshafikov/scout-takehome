package helpers

import (
	"crypto/rand"
	"math/big"
	"time"
)

func StrStartsWith(search, str string) bool {
	return len(str) >= len(search) && str[:len(search)] == search
}

func GetPointer[T string | int | uint | int64 | uint64 | bool | time.Time](val T) *T {
	return &val
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateRandomString(length int) (string, error) {
	b := make([]rune, length)
	for i := range b {
		randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}

		b[i] = letters[randomNumber.Int64()]
	}

	return string(b), nil
}

func GenerateRandomNumberInRange(min, max int) (int, error) {
	random, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(random.Int64()) + min, nil
}
