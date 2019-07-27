package store

import (
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
	values := v.store.ReadAtPos(v.pos)
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

//Store is the defintion for the storage implementation requirements
type Store interface {
	Get(id string) (Vector, error)
	Set(id string, vector Vector) error
	Delete(id string) error
	KNN(vector Vector, k int) (*[]Distance, error)
}
