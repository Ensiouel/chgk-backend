package main

import (
	"math/rand"
	"time"
)

func ranges(n int) (s []int) {
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return
}

func randPop(s []int) (int, []int) {
	if len(s) == 1 {
		return s[0], []int{}
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(s))
	v := s[n]

	return v, append(s[:n], s[n+1:]...)
}
