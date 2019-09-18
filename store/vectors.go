package store

import (
	"math"
	"math/rand"
)

type Vector interface {
	Values() *[]float64
	Name() string
}

//Vector is the main storage unit
type MemoryVector struct {
	ID    string
	Array []float64
}

func (v *MemoryVector) Values() *[]float64 {
	return &v.Array
}
func (v *MemoryVector) Name() string {
	return v.ID
}

type PersistantVector struct {
	ID    string
	pos   uint32
	store *PersistantStore
}

func (v *PersistantVector) Values() *[]float64 {
	bytes, _ := v.store.vectorsFile.Get([]byte(v.ID), nil)
	values := BytesToArray(bytes)
	return &values
}
func (v *PersistantVector) Name() string {
	return v.ID
}

//Zeros will return a Vector where all the values are equal to 0.0
func Zeros(id string, size int) *MemoryVector {
	return &MemoryVector{id, make([]float64, size)}
}

//Ones will return a Vector where all the values are equal to 1.0
func Ones(id string, size int) *MemoryVector {
	var v = MemoryVector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Array[q] = 1
	}
	return &v
}

//Random will create a Vector with values in range (0,1)
func Random(id string, size int) *MemoryVector {
	var v = MemoryVector{id, make([]float64, size)}
	for q := 0; q < size; q++ {
		v.Array[q] = rand.Float64()
	}
	return &v
}

/**
	* EuclideanDistance
 	* Distance between two vectors x= <x_1, x_2, x_3> and y= <y_1, y_2, y_3>
 	* is defined as {(x_1-y_1)²+(x_2-y_2)²+(x_3-y_3)²\}^{1/2}.
*/
func EuclideanDistance(x Vector, y Vector) float64 {
	var dist float64

	i := 0
	for i < len(*x.Values()) {
		left := (*x.Values())[i]
		right := (*y.Values())[i]
		dist += math.Pow((left - right), 2)
		i++
	}
	return math.Sqrt(dist)
}

/**
* Return structure for the kNN algorithm
* holding the string as key, the vector for corresponding key and the distance
 */
type Distance struct {
	Target   Vector
	Distance float64 //Distance
}
