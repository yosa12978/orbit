package pkg

import (
	"fmt"
	"math"
	"sort"
)

var charset string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(n int64) string {
	id := ""
	for n > 0 {
		id = string(charset[n%int64(len(charset))]) + id
		n /= int64(len(charset))
	}
	return id
}

func DecodeBase62(s string) (int64, error) {
	var idInt int64 = 0
	for k, v := range s {
		index := sort.Search(len(charset), func(i int) bool {
			return charset[i] >= byte(v)
		})
		if index < len(charset) && charset[index] == byte(v) {
			idInt += int64(index) * int64(math.Pow(float64(len(charset)), float64(len(s)-k-1)))
			continue
		}
		return 0, fmt.Errorf("'%c' not found in the alphabet\n", v)
	}
	return idInt, nil
}
