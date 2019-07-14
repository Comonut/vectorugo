package store

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	//train vectors
	var trainVector1 = Ones("v1", 5)
	var trainVector2 = Zeros("v2", 5)
	var trainVector3 = Random("v3", 5)
	var trainVector4 = Random("v4", 5)
	//test vector
	var testVector = Random("t1", 5)

	vectorMap := map[string]*Vector{
		"v1": trainVector1,
		"v2": trainVector2,
		"v3": trainVector3,
		"v4": trainVector4,
	}

	neighbors := KNN(vectorMap, testVector, 3)
	for i := range neighbors{
		fmt.Printf("Neighbors: %v", neighbors[i])
		fmt.Println()
	}
}
