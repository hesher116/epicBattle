package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rnd  *rand.Rand
	once sync.Once
)

func initRandom() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int) int {
	once.Do(initRandom)
	return min + rnd.Intn(max-min+1)
}
