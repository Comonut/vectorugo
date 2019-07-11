package store

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	//train vectors
	var trainVector1 = Ones("", 32)
	var trainVector2 = Zeros("", 32)
	var trainVector3 = Random("", 32)
	var trainVector4 = Random("", 32)
	//test vector
	var testVector = Random("", 32)

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

