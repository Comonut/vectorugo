package store

import (
	"math"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

type PersistantStore struct {
	dimension   uint32 //size of vectors
	size        uint32 //number of vectors inside
	keys        map[string]bool
	seachIndex  *Index
	vectorsFile *leveldb.DB //vectors file (see encodings.go)
}

func NewPersitantStore(dimension uint32, indexFile, vectorsFile, searchindexFile string) (*PersistantStore, error) {
	_, err1 := os.Stat(vectorsFile)
	// If index doesn't exist
	if os.IsNotExist(err1) {
		return ConstructPersistantStore(dimension, indexFile, vectorsFile, searchindexFile), nil
	} else { //if they don't exist create a new storage
		return LoadPersistantStore(dimension, indexFile, vectorsFile, searchindexFile), nil
	}
}

//Creates a new persistant store
func ConstructPersistantStore(dimension uint32, indexFile, vectorsFile, searchindexFile string) *PersistantStore {
	logrus.Info("New persistant storage initialized")
	vectors, _ := leveldb.OpenFile(vectorsFile, nil)
	store := &PersistantStore{
		dimension:   dimension,
		size:        uint32(0),
		vectorsFile: vectors,
		keys:        make(map[string]bool),
		seachIndex:  nil}

	store.seachIndex = NewIndex(vectors, store)
	return store
}

//Loads an existing store
func LoadPersistantStore(dimension uint32, indexFile, vectorsFile, searchindexFiles string) *PersistantStore {
	logrus.Info("Loading existant persistance storage")
	vectors, _ := leveldb.OpenFile(vectorsFile, nil)
	iter := vectors.NewIterator(nil, nil)
	counter := uint32(0)
	keys := make(map[string]bool)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		keys[string(iter.Key())] = true
		counter++
	}
	iter.Release()
	logrus.Info("Done!") //for each line in the id's list - map it to it's position
	store := &PersistantStore{
		dimension:   dimension,
		size:        counter,
		keys:        keys,
		seachIndex:  nil,
		vectorsFile: vectors}

	store.seachIndex = LoadIndex(vectors, store)
	return store
}

func (s *PersistantStore) Set(id string, vector Vector) error {

	persistant := &PersistantVector{ID: vector.Name(), store: s}
	if _, ok := s.keys[id]; ok { //if this id is indexed
		s.WriteVector(vector.Values(), vector.Name())
	} else {
		s.WriteVector(vector.Values(), vector.Name())
		s.keys[id] = true
		s.seachIndex.maxlen = 2 * int(math.Sqrt(float64(s.size)))
		s.seachIndex.AddVector(persistant)
		s.size++ //increment size of store

	}
	return nil
}

func (s *PersistantStore) Get(id string) (Vector, error) {

	values, err := s.ReadVector(id)

	if err != nil {
		return nil, err
	}
	return &MemoryVector{
		ID:    id,
		Array: values}, nil
}

func (s *PersistantStore) Delete(id string) error {
	return s.vectorsFile.Delete([]byte(id), nil)
}

func (s *PersistantStore) KNN(vector Vector, k int) (*[]Distance, error) {
	return s.seachIndex.IndexKNN(k, vector)
}
