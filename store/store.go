package store

import (
	"math/rand"
)

type Vector struct {
	Id     string
	Values []float64
}

func Zeros(id string, size int) *Vector {
	return &Vector{id, make([]float64, size)}
}

func Ones(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Values[q] = 1
	}
	return &v
}

func Random(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Values[q] = rand.Float64()
	}
	return &v
}

func (v *Vector) Sum() float64 {
	var sum = 0.0
	for _, element := range v.Values {
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
