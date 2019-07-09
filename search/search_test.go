package search

import (
	"fmt"
	"testing"
	"github.com/Comonut/vectorugo/store"
)

func TestSearch(t *testing.T) {
	//train vectors
	var trainVector1 = store.Ones("", 32)
	var trainVector2 = store.Zeros("", 32)
	var trainVector3 = store.Random("", 32)
	var trainVector4 = store.Random("", 32)
	//test vector
	var testVector = store.Random("", 32)

	type Vector = store.Vector
	vectorMap := map[string]*Vector{
		"v1": trainVector1,
		"v2": trainVector2,
		"v3": trainVector3,
		"v4": trainVector4,
	}

	a := LabelVectors(vectorMap, testVector)
	fmt.Printf("Distances: %v", a)
	fmt.Println()

	b := GetNeigbors(a, 2)
	fmt.Printf("Neighbors: %v", b)
	fmt.Println()
}
