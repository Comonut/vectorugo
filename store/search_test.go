package store

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	//train vectors
	var v1 = Zeros("v1", 5)
	var v2 = Zeros("v2", 5)
	var v3 = Zeros("v3", 5)
	v3.Values[3] = 1
	var v4 = Zeros("v4", 5)
	v4.Values[3] = 5
	var v5 = Zeros("v5", 5)
	v5.Values[3] = 2
	v5.Values[4] = 1

	s := NewSimpleMapStore()
	s.Set(v5.ID, v5)
	s.Set(v3.ID, v3)
	s.Set(v2.ID, v2)
	s.Set(v4.ID, v4)
	s.Set(v1.ID, v1)

	neighbors := MapStoreKNN(&s, v1, 5)
	for i := range *neighbors {
		fmt.Printf("Neighbors: %v", (*neighbors)[i])
		fmt.Println()
	}
}
