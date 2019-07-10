package store

import (
	"errors"
	"math/rand"
)

type Vector struct {
	id     string
	values []float64
}

func Zeros(id string, size int) *Vector {
	return &Vector{id, make([]float64, size)}
}

func Ones(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.values[q] = 1
	}
	return &v
}

func Random(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.values[q] = rand.Float64()
	}
	return &v
}

func (v *Vector) Sum() float64 {
	var sum = 0.0
	for _, element := range v.values {
		sum += element
	}
	return sum
}

type Store interface {
	Get(id string) (*Vector, error)
	Set(id string, vector *Vector) error
	Delete(id string) error
	KNN(vector *Vector, k int) ([]*Vector, error)
}

type SimpleMapStore map[string]*Vector

func newSimpleMapStore() SimpleMapStore {
	return make(map[string]*Vector)
}

func (s *SimpleMapStore) Get(id string) (*Vector, error) {
	var value, found = (*s)[id]
	if found {
		return value, nil
	}
	return value, errors.New("id not found")
}

func (s *SimpleMapStore) Set(id string, vector *Vector) error {
	(*s)[id] = vector
	return nil
}

func (s *SimpleMapStore) Delete(id string) error {
	delete(*s, id)
	return nil
}

func (s *SimpleMapStore) KNN(vector *Vector, k int) ([]*Vector, error) {
	return nil, nil
}
