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
	* Return structure for the kNN algorithm
	* holding the string as key, the vector for corresponding key and the distance
 */
type Distance struct {
	Key      string    //Vector ID
	Values   []float64 //Vector Values
	Distance float64   //Distance
}

/**
	* KNN
	* take vector map, test vector and number of neighbors
	* return array with closest k neighbors containing Id, Vectro Values and Distance for each neighbor
 */
func KNN(storeVectors map[string]*Vector, testVector *Vector, k int) []Distance {
	sortedDistances := make([]Distance, len(storeVectors))
	var currentVector Distance
	counter := 0

	//loop through map, calculate distance for each vector, append result in return array
	for _, v := range storeVectors{
		currentVector.Key = v.Id
		currentVector.Values = v.Values
		currentVector.Distance = EuclideanDistance(v, testVector)
		sortedDistances[counter] = currentVector
		counter++
	}

	//sort array by distance in incrementing order
	sort.Slice(sortedDistances, func(i, j int) bool {
		return sortedDistances[i].Distance < sortedDistances[j].Distance
	})

	return sortedDistances[0:k]
}
