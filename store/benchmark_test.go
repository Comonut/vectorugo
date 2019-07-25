package store

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

// STORAGE BENCHMARKS
// YOU CAN RUN WITH 'go test ./store -bench=.'
// TODO: SEARCH BENCHMARKS

//FOR PREVENTING COMPILER OPTIMIZATIONS
var e error
var f *Vector

func BenchmarkInsertMemoryStore(b *testing.B) {
	store := NewSimpleMapStore()
	var err error
	for n := 0; n < b.N; n++ {
		err = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	e = err
}

func BenchmarkGetMemoryStore(b *testing.B) {
	store := NewSimpleMapStore()
	for n := 0; n < 100; n++ {
		_ = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	var err error
	for n := 0; n < b.N; n++ {
		f, err = store.Get(strconv.Itoa(rand.Intn(100)))
	}
	e = err
}

func BenchmarkInsertPersistantStore(b *testing.B) {
	store, _ := NewPersitantStore(256, "benchmark_index.bin", "benchmark_vectors.bin")
	var err error
	for n := 0; n < b.N; n++ {
		err = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	e = err
	os.Remove("benchmark_index.bin")
	os.Remove("benchmark_vectors.bin")
}

func BenchmarkGetPersistantStore(b *testing.B) {
	store, _ := NewPersitantStore(256, "benchmark_index.bin", "benchmark_vectors.bin")
	for n := 0; n < 100; n++ {
		_ = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	var err error
	var v *Vector
	for n := 0; n < b.N; n++ {
		v, err = store.Get(strconv.Itoa(rand.Intn(100)))
	}
	e = err
	f = v
	os.Remove("benchmark_index.bin")
	os.Remove("benchmark_vectors.bin")
}
