package store

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

type PersistantStore struct {
	dimension   uint32            //size of vectors
	index       map[string]uint32 //maps id's to positions in the file
	size        uint32            //number of vectors inside
	seachIndex  *Index
	indexFile   *os.File    //file that contains the mapping from id to position in vectors file
	vectorsFile *leveldb.DB //vectors file (see encodings.go)
}

func NewPersitantStore(dimension uint32, indexFile, vectorsFile, searchindexFile string) (*PersistantStore, error) {
	_, err1 := os.Stat(indexFile)
	_, err2 := os.Stat(indexFile)

	// If index doesn't exist
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		return ConstructPersistantStore(dimension, indexFile, vectorsFile, searchindexFile), nil
	} else if os.IsNotExist(err1) || os.IsNotExist(err2) { //If one of the two files exists - something is wrong
		return nil, fmt.Errorf("Index or vectors file exists, but other one doesn't")
	} else { //if they don't exist create a new storage
		return LoadPersistantStore(dimension, indexFile, vectorsFile, searchindexFile), nil
	}
}

//Creates a new persistant store
func ConstructPersistantStore(dimension uint32, indexFile, vectorsFile, searchindexFile string) *PersistantStore {
	logrus.Info("New persistant storage initialized")
	index, _ := os.Create(indexFile)
	vectors, _ := leveldb.OpenFile(vectorsFile, nil)
	store := &PersistantStore{
		dimension:   dimension,
		index:       make(map[string]uint32),
		size:        uint32(0),
		indexFile:   index,
		vectorsFile: vectors,
		seachIndex:  nil}

	store.seachIndex = NewIndex(vectors, store)
	return store
}

//Loads an existing store
func LoadPersistantStore(dimension uint32, indexFile, vectorsFile, searchindexFiles string) *PersistantStore {
	logrus.Info("Loading existant persistance storage")
	index, _ := os.OpenFile(indexFile, os.O_RDWR|os.O_CREATE, 0755)
	vectors, _ := leveldb.OpenFile(vectorsFile, nil)
	scanner := bufio.NewScanner(index)
	posIndex := make(map[string]uint32)        //id to pos index
	inversePosIndex := make(map[uint32]string) //pos to id index
	counter := uint32(0)
	for scanner.Scan() {
		name := scanner.Text()
		posIndex[name] = counter
		inversePosIndex[counter] = name
		counter++
	}

	logrus.Info("Done!") //for each line in the id's list - map it to it's position
	store := &PersistantStore{
		dimension:   dimension,
		index:       posIndex,
		size:        counter,
		indexFile:   index,
		seachIndex:  nil,
		vectorsFile: vectors}

	store.seachIndex = LoadIndex(vectors, store)
	return store
}

func (s *PersistantStore) Set(id string, vector Vector) error {

	persistant := &PersistantVector{ID: vector.Name(), pos: 0, store: s}
	if pos, ok := s.index[id]; ok { //if this id is indexed
		persistant.pos = pos
		s.WriteVector(vector.Values(), vector.Name())
	} else {
		persistant.pos = s.size
		s.WriteVector(vector.Values(), vector.Name())
		_, err := s.indexFile.WriteString(id + "\n")
		if err != nil {
			return err
		}
		s.index[id] = s.size //create pos index for id
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
