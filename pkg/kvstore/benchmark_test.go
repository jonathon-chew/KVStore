package kvstore

import (
	"strconv"
	"testing"
)

func benchSet(i int, s Store, b *testing.B) {

	words, err := findCollisions(10, i)
	if err != nil {
		b.Fatalf("[ERROR]: can't make the word list")
	}

	newWords := []string{}
	for idx, wd := range words {
		newWord := wd + strconv.Itoa(idx)
		newWords = append(newWords, newWord)
	}

	for i := 0; i < b.N; i++ {
		for _, wf := range newWords {
			s.Add(wf, reverse(wf))
		}
	}
}

// func Benchmark_SetBench(b *testing.B) {

// 	s := NewKVStore(10)

// 	for i := 0; i < b.N; i++ {
// 		benchSet(10, s, b)
// 		benchSet(20, s, b)
// 		benchSet(40, s, b)
// 	}
// }

func BenchmarkSet1(b *testing.B) {
	s := NewKVStore(10)
	benchSet(1, s, b)
}
func BenchmarkSet2(b *testing.B) {
	s := NewKVStore(10)
	benchSet(2, s, b)
}
func BenchmarkSet3(b *testing.B) {
	s := NewKVStore(10)
	benchSet(3, s, b)
}
func BenchmarkSet10(b *testing.B) {
	s := NewKVStore(10)
	benchSet(10, s, b)
}
func BenchmarkSet20(b *testing.B) {
	s := NewKVStore(10)
	benchSet(20, s, b)
}
func BenchmarkSet40(b *testing.B) {
	s := NewKVStore(10)
	benchSet(40, s, b)
}
