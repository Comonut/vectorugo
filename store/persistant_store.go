package store

import (
	"bufio"
	"fmt"
	"os"
)

type PersistantStore struct {
	dimension   uint32
	index       map[string]uint32
	size        uint32
	indexFile   *os.File
	vectorsFile *os.File
}

func NewPersitantStore(dimension uint32, indexFile, vectorsFile string) (*PersistantStore, error) {
	_, err1 := os.Stat(indexFile)
	_, err2 := os.Stat(indexFile)

	// create file if not exists
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		return ConstructPersistantStore(dimension, indexFile, vectorsFile), nil
	} else if os.IsNotExist(err1) || os.IsNotExist(err2) {
		return nil, fmt.Errorf("Index or vectors file exists, but other one doesn't")
	} else {
		return LoadPersistantStore(dimension, indexFile, vectorsFile), nil
	}
}

func ConstructPersistantStore(dimension uint32, indexFile, vectorsFile string) *PersistantStore {
	fmt.Println("New persistant storage initialized")
	index, _ := os.Create(indexFile)
	vectors, _ := os.Create(vectorsFile)
	return &PersistantStore{
		dimension:   dimension,
		index:       make(map[string]uint32),
		size:        uint32(0),
		indexFile:   index,
		vectorsFile: vectors}
}

func LoadPersistantStore(dimension uint32, indexFile, vectorsFile string) *PersistantStore {
	fmt.Println("Loading existant persistance storage")
	index, _ := os.OpenFile(indexFile, os.O_RDWR|os.O_CREATE, 0755)
	vectors, _ := os.OpenFile(vectorsFile, os.O_RDWR|os.O_CREATE, 0755)
	scanner := bufio.NewScanner(index)
	posIndex := make(map[string]uint32)
	counter := uint32(0)
	fmt.Println("Done!")
	for scanner.Scan() {
		posIndex[scanner.Text()] = counter
		counter++
	}
	return &PersistantStore{
		dimension:   dimension,
		index:       posIndex,
		size:        counter,
		indexFile:   index,
		vectorsFile: vectors}
}

func (s *PersistantStore) Set(id string, vector *Vector) error {
	if pos, ok := s.index[id]; ok {
		s.WriteAtPos(vector, pos)
	} else {
		s.WriteAtPos(vector, s.size)
		s.indexFile.WriteString(id + "\n")
		s.index[id] = s.size
		s.size++
	}
	return nil
}

func (s *PersistantStore) Get(id string) (*Vector, error) {
	pos, ok := s.index[id]
	if !ok {
		return nil, fmt.Errorf("value not present")
	}

	return &Vector{
		ID:     id,
		Values: s.ReadAtPos(pos)}, nil
}

func (s *PersistantStore) Delete(id string) error {
	return nil
}

func (s *PersistantStore) KNN(vector *Vector, k int) (*[]Distance, error) {
	return nil, nil
}
