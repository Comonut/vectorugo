package store

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestIndex(t *testing.T) {
	i := NewIndex()

	var v *Vector

	for q := 0; q < 1000; q++ {
		if rand.Intn(2)%2 == 0 {
			v = Zeros("000-"+strconv.Itoa(q), 256)
		} else {
			v = Ones("111-"+strconv.Itoa(q), 256)
		}
		i.AddVector(v)
		i.maxlen = 2 * int(math.Sqrt(float64(q+1)))
	}

	distances := i.IndexKNN(5, v)
	fmt.Print(distances)
}
