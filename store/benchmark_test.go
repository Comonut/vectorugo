package store

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

// STORAGE BENCHMARKS
// YOU CAN RUN WITH 'go test ./store -bench=.'

//FOR PREVENTING COMPILER OPTIMIZATIONS
var e error
var f Vector
var d *[]Distance

func BenchmarkInsertMemoryStore(b *testing.B) {
	store := NewSimpleMapStore()
	var err error
	b.ResetTimer()
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
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f, err = store.Get(strconv.Itoa(rand.Intn(100)))
	}
	e = err
}

func BenchmarkInsertPersistantStore(b *testing.B) {
	store, _ := NewPersitantStore(256, "benchmark_index.bin", "benchmark_vectors.bin", "benchmark_search.bin")
	var err error
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	e = err
	os.Remove("benchmark_index.bin")
	os.RemoveAll("benchmark_vectors.bin")
	os.Remove("benchmark_search.bin")
}

func BenchmarkGetPersistantStore(b *testing.B) {
	store, _ := NewPersitantStore(256, "benchmark_index.bin", "benchmark_vectors.bin", "benchmark_search.bin")
	for n := 0; n < 100; n++ {
		_ = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	var err error
	var v Vector
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v, err = store.Get(strconv.Itoa(rand.Intn(100)))
	}
	e = err
	f = v
	os.Remove("benchmark_index.bin")
	os.RemoveAll("benchmark_vectors.bin")
	os.Remove("benchmark_search.bin")
}

func BenchmarkKNNsearch(b *testing.B) {
	store, _ := NewPersitantStore(256, "benchmark_index.bin", "benchmark_vectors.bin", "benchmark_search.bin")
	for n := 0; n < 100; n++ {
		_ = store.Set(strconv.Itoa(n), Random(strconv.Itoa(n), 256))
	}
	var err error
	var d2 *[]Distance
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d2, err = store.KNN(Random("testV", 256), 5)
	}
	e = err
	d = d2
	os.Remove("benchmark_index.bin")
	os.RemoveAll("benchmark_vectors.bin")
	os.Remove("benchmark_search.bin")
}
