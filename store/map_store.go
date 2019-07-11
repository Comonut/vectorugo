package store

import "errors"

type SimpleMapStore map[string]*Vector

func newSimpleMapStore() SimpleMapStore {
	return make(map[string]*Vector)
}

func (s *SimpleMapStore) Get(id string) (*Vector, error) {
	var value, found = (*s)[id]
	if found {
		return value, nil
	}
	return value, errors.New("id not found")
}

func (s *SimpleMapStore) Set(id string, vector *Vector) error {
	(*s)[id] = vector
	return nil
}

func (s *SimpleMapStore) Delete(id string) error {
	delete(*s, id)
	return nil
}

func (s *SimpleMapStore) KNN(vector *Vector, k int) ([]*Vector, error) {
	return nil, nil
}
