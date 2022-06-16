package codekit

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
	"sync"
	"time"
)

type randomizer struct {
	sync.Mutex
	r *rand.Rand
}

func (r *randomizer) intN(limit int) int {
	defer r.Unlock()
	r.Lock()
	return r.r.Intn(limit)
}

var (
	once sync.Once
	r    *randomizer
)

func initRandomSource() {
	once.Do(func() {
		src := rand.NewSource(time.Now().UnixNano())
		r = new(randomizer)
		r.r = rand.New(src)
	})
}

func RandInt(limit int) int {
	initRandomSource()
	return r.intN(limit)
}

func RandFloat(limit int, decimal int) float64 {
	flim := float64(limit)
	fdec := float64(decimal)
	initRandomSource()
	powerLimit := int(flim * math.Pow(10, fdec))
	randPower := r.intN(powerLimit)
	return float64(randPower) / math.Pow(10, fdec)
}

func MD5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(baseChars string, n int) string {
	if baseChars == "" {
		baseChars = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz_!"
	}
	baseCharsLen := len(baseChars)

	rnd := ""
	for i := 0; i < n; i++ {
		x := RandInt(baseCharsLen)
		rnd += string(baseChars[x])
	}
	return rnd
}

func RandomString(length int) string {
	return GenerateRandomString("", length)
}
