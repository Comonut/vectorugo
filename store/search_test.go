package store

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	//train vectors
	var trainVector1 = Ones("", 5)
	var trainVector2 = Zeros("", 5)
	var trainVector3 = Random("", 5)
	var trainVector4 = Random("", 5)
	//test vector
	var testVector = Random("", 5)

	vectorMap := map[string]*Vector{
		"v1": trainVector1,
		"v2": trainVector2,
		"v3": trainVector3,
		"v4": trainVector4,
	}

	a := LabelVectors(vectorMap, testVector)
	fmt.Printf("Distances: %v", a)
	fmt.Println()

	b := GetNeigbors(vectorMap, a, 2)
	fmt.Printf("Test vector: %v %v", testVector.Id, testVector.Values)
	fmt.Println()

	for k, v := range b{
		fmt.Printf("Neighbors: %v %v", k, v)
		fmt.Println()
	}
}
