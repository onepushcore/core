package core

import (
	"math/rand"
	"time"
)

const Base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, length)
	for i := range res {
		res[i] = Base62Chars[rnd.Intn(len(Base62Chars))]
	}
	return string(res)
}

func ToBase62(number int) string {
	result := ""
	for number > 0 {
		result = string(Base62Chars[number%62]) + result
		number /= 62
	}
	return result
}
