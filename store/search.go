package store

import (
	"math"
	"sort"
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
func LabelVectors(m SimpleMapStore, test *Vector) []float64 {
	a := make([]float64, len(m))
	i := 0

	for _, v := range m {
		a[i] = EuclideanDistance(v, test)
		i++
	}

	return a
}

/**
* GetNeigbors
* Sort the distance and determine nearest neighbors based on the k-th minimum distance.
 */
func GetNeigbors(a []float64, k int) []float64 {

	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	return a[0:k]
}
