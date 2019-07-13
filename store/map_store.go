package store

import "errors"

//SimpleMapStore  implements the Store interface using a Go map
type SimpleMapStore map[string]*Vector

func NewSimpleMapStore() SimpleMapStore {
	return make(map[string]*Vector)
}

//Get Vector pointer using a given id, returns error if not found
func (s *SimpleMapStore) Get(id string) (*Vector, error) {
	var value, found = (*s)[id]
	if found {
		return value, nil
	}
	return value, errors.New("id not found")
}

//Set Vector for a given id
func (s *SimpleMapStore) Set(id string, vector *Vector) error {
	(*s)[id] = vector
	return nil
}

//Delete a given value from the store
func (s *SimpleMapStore) Delete(id string) error {
	delete(*s, id)
	return nil
}

//KNN returns the K nearest neighbours to a given vector
//This vector does not have to be inside the store.
func (s *SimpleMapStore) KNN(vector *Vector, k int) ([]*Vector, error) {
	return nil, nil
}
