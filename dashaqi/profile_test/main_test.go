package main

import (
	"github.com/songrenru/dashaqi"
	"testing"
)

func TestAdd(t *testing.T) {
	_ = dashaqi.Add("go-programming-tour-book")
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dashaqi.Add("go-programming-tour-book")
	}
}
