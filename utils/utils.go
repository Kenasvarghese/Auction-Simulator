package utils

import (
	"fmt"
	"math/rand"
)

func MakeAttributes() map[string]string {
	m := make(map[string]string, 20)
	for i := range 20 {
		m[fmt.Sprintf("attr_%02d", i+1)] = fmt.Sprintf("val_%d", rand.Intn(10000))
	}
	return m
}
