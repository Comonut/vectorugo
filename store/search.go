package store

import (
	"math"
	"sort"
	"strings"
)

/**
	* EuclideanDistance
 	* Distance between two vectors x= <x_1, x_2, x_3> and y= <y_1, y_2, y_3>
 	* is defined as {(x_1-y_1)²+(x_2-y_2)²+(x_3-y_3)²\}^{1/2}.
*/
func EuclideanDistance(x *Vector, y *Vector) float64 {
	var dist float64

	i := 0
	for i < len(x.Values) {
		left := x.Values[i]
		right := y.Values[i]
		dist += math.Pow((left - right), 2)
		i++
	}
	return math.Sqrt(dist)
}

/**
* LabelVectors
* Iterate through map of vectors, determine Euclidean Distance between current vector and
* test vector, return array with distances
 */
func LabelVectors(m SimpleMapStore, test *Vector) map[string]float64 {
	a := make(map[string]float64)

	for i, v := range m {
		a[i] = EuclideanDistance(v, test)
	}

	return a
}

/**
* Return structure for the kNN algorithm
* holding the string as key and distance as value
 */
type Kv struct {
	Key   string  //string ID
	Value float64 //vector distance
}

/**
* GetNeigbors
* Sort the distance and determine nearest neighbors based on the k-th minimum distance.
 */
func GetNeigbors(m map[string]*Vector, a map[string]float64, k int) map[string]*Vector {
	var sortedStructs []Kv
	finalResult := make(map[string]*Vector)

	//array with (id, distance) pairs
	for k, v := range a {
		sortedStructs = append(sortedStructs, Kv{k, v})
	}

	//sort array by distance in incrementing order
	sort.Slice(sortedStructs, func(i, j int) bool {
		return sortedStructs[i].Value < sortedStructs[j].Value
	})

	//interate over map with vectors
	for key, v := range m {
		//for each entry check if it exists in the sorted array up to the number of neighbors we look for
		for j := range sortedStructs[0:k] {
			//if entry found, add it to the finalResult map and exit current iteration
			if strings.Compare(key, sortedStructs[j].Key) == 0 {
				//jackpot
				finalResult[sortedStructs[j].Key] = v
				break
			}
		}
	}
	return finalResult
}
