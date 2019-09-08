package store

import (
	"errors"
	"math"
)

//SimpleMapStore  implements the Store interface using a Go map
type SimpleMapStore struct {
	vectors map[string]Vector
	index   *Index
}

//NewSimpleMapStore returns a new map store
func NewSimpleMapStore() *SimpleMapStore {
	return &SimpleMapStore{vectors: make(map[string]Vector), index: NewIndex(nil, nil)}
}

//Get Vector pointer using a given id, returns error if not found
func (s *SimpleMapStore) Get(id string) (Vector, error) {
	var value, found = s.vectors[id]
	if found {
		return value, nil
	}
	return value, errors.New("id not found")
}

//Set Vector for a given id
func (s *SimpleMapStore) Set(id string, vector Vector) error {
	s.vectors[id] = vector
	s.index.AddVector(vector)
	s.index.maxlen = 2 * int(math.Sqrt(float64(len(s.vectors))))
	return nil
}

//Delete a given value from the store
func (s *SimpleMapStore) Delete(id string) error {
	delete(s.vectors, id)
	return nil
}

//KNN returns the K nearest neighbours to a given vector
//This vector does not have to be inside the store.
func (s *SimpleMapStore) KNN(vector Vector, k int) (*[]Distance, error) {
	return s.index.IndexKNN(k, vector)
}
