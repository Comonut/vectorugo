package store

import (
	"math/rand"
)

//Vector is the main storage unit
type Vector struct {
	ID string

	Values []float64
}

//Zeros will return a Vector where all the values are equal to 0.0
func Zeros(id string, size int) *Vector {
	return &Vector{id, make([]float64, size)}
}

//Ones will return a Vector where all the values are equal to 1.0
func Ones(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Values[q] = 1
	}
	return &v
}

//Random will create a Vector with values in range (0,1)
func Random(id string, size int) *Vector {
	var v = Vector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Values[q] = rand.Float64()
	}
	return &v
}

//Sum of the vector elements
func (v *Vector) Sum() float64 {
	var sum = 0.0
	for _, element := range v.Values {
		sum += element
	}
	return sum
}

//Store is the defintion for the storage implementation requirements
type Store interface {
	Get(id string) (*Vector, error)
	Set(id string, vector *Vector) error
	Delete(id string) error
	KNN(vector *Vector, k int) (*[]Distance, error)
}
