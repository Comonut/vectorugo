package store

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type PersistantStore struct {
	dimension   uint32            //size of vectors
	index       map[string]uint32 //maps id's to positions in the file
	size        uint32            //number of vectors inside
	indexFile   *os.File          //file that contains the mapping from id to position in vectors file
	vectorsFile *os.File          //vectors file (see encodings.go)
}

func NewPersitantStore(dimension uint32, indexFile, vectorsFile string) (*PersistantStore, error) {
	_, err1 := os.Stat(indexFile)
	_, err2 := os.Stat(indexFile)

	// If index doesn't exist
	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		return ConstructPersistantStore(dimension, indexFile, vectorsFile), nil
	} else if os.IsNotExist(err1) || os.IsNotExist(err2) { //If one of the two files exists - something is wrong
		return nil, fmt.Errorf("Index or vectors file exists, but other one doesn't")
	} else { //if they don't exist create a new storage
		return LoadPersistantStore(dimension, indexFile, vectorsFile), nil
	}
}

//Creates a new persistant store
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

//Loads an existing store
func LoadPersistantStore(dimension uint32, indexFile, vectorsFile string) *PersistantStore {
	fmt.Println("Loading existant persistance storage")
	index, _ := os.OpenFile(indexFile, os.O_RDWR|os.O_CREATE, 0755)
	vectors, _ := os.OpenFile(vectorsFile, os.O_RDWR|os.O_CREATE, 0755)
	scanner := bufio.NewScanner(index)
	posIndex := make(map[string]uint32) //id to pos index
	counter := uint32(0)
	for scanner.Scan() {
		posIndex[scanner.Text()] = counter
		counter++
	}
	fmt.Println("Done!") //for each line in the id's list - map it to it's position
	return &PersistantStore{
		dimension:   dimension,
		index:       posIndex,
		size:        counter,
		indexFile:   index,
		vectorsFile: vectors}
}

func (s *PersistantStore) Set(id string, vector *Vector) error {
	if pos, ok := s.index[id]; ok { //if this id is indexed
		s.WriteAtPos(vector, pos) //write at vector's position in file
	} else {
		s.WriteAtPos(vector, s.size)                 //if it doesn't exist, write at end of vector's file
		_, err := s.indexFile.WriteString(id + "\n") //add it to list of ids
		if err != nil {
			return err
		}
		s.index[id] = s.size //create pos index for id
		s.size++             //increment size of store
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

type idPosPair struct {
	id  string
	pos uint32
}

func (s *PersistantStore) Delete(id string) error {
	positionsArr := make([]idPosPair, len(s.index))
	counter := 0
	for k, v := range s.index {
		positionsArr[counter] = idPosPair{id: k, pos: v}
		counter++
	}
	sort.Slice(positionsArr, func(i, j int) bool {
		return positionsArr[i].pos < positionsArr[j].pos
	})
	last := positionsArr[len(positionsArr)-1]
	positionsArr[s.index[id]] = last
	s.index[last.id] = s.index[id]
	delete(s.index, id)
	s.WriteAtPos(&Vector{"", s.ReadAtPos(last.pos)}, s.index[last.id])
	s.size--
	err := s.vectorsFile.Truncate(int64(s.size * s.dimension * 8))
	if err != nil {
		panic(fmt.Errorf("error shrinking vectors file"))
	}
	index := s.indexFile.Name()
	os.Remove(index)
	s.indexFile, _ = os.Create(index)
	for _, k := range positionsArr[:s.size] {
		_, err = s.indexFile.WriteString(k.id + "\n")
		if err != nil {
			panic(fmt.Errorf("error updanting index file"))
		}
	}
	return nil
}

func (s *PersistantStore) KNN(vector *Vector, k int) (*[]Distance, error) {
	return nil, nil
}