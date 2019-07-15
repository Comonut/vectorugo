package store

import "os"

type PersistantStore struct {
	dimension   int
	index       map[string]uint32
	size        uint32
	indexFile   *os.File
	vectorsFile *os.File
}

func NewPersitantStore(dimension int, indexFile, vectorsFile string) *PersistantStore {
	index, _ := os.Create(indexFile)
	vectors, _ := os.Create(vectorsFile)
	return &PersistantStore{
		dimension:   dimension,
		index:       make(map[string]uint32),
		size:        uint32(0),
		indexFile:   index,
		vectorsFile: vectors}
}
