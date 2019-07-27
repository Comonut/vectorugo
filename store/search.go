package store

import (
	"math"
)

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
