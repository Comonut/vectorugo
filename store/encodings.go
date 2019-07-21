package store

import (
	"encoding/binary"
	"fmt"
	"math"
)

/*
For handling writing of vectors into a file
The file consists of float64 bytes, so for example if you have one float you will have 8 bytes
if you have an array of size D - you will have 8*D bytes
if you have N array's of size D - you will have N * 8*D bytes
To read or write a specific vector, you use his pos - which is his numeric index
*/

//Transforms a float64 into 8 bytes
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

//Transforms 8 bytes into a float
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

//transforms N floats into N*8 bytes
func ArrToBytes(arr []float64) []byte {
	bytes := make([]byte, 8*len(arr))
	var part []byte
	for i, f := range arr {
		part = Float64bytes(f)
		copy(bytes[i*8:i*8+8], part)
	}
	return bytes
}

//Transforms N*8 bytes into N floats
func BytesToArray(arr []byte) []float64 {
	if len(arr)%8 != 0 {
		panic(fmt.Errorf("corrupted vector binaries"))
	}
	floats := make([]float64, len(arr)/8)
	for i := 0; i < len(arr)/8; i++ {
		floats[i] = Float64frombytes(arr[i*8 : i*8+8])

	}
	return floats
}

//Writes the values of a vector in a given position in the file
func (s *PersistantStore) WriteAtPos(v *Vector, pos uint32) {
	bytes := ArrToBytes(v.Values)
	_, err := s.vectorsFile.WriteAt(bytes, int64(pos*8*s.dimension))
	if err != nil {
		panic(fmt.Errorf("Error writing in vectors file"))
	}
}

//Reads the values of a vector from a given position
func (s *PersistantStore) ReadAtPos(pos uint32) []float64 {
	bytes := make([]byte, s.dimension*8)
	_, err := s.vectorsFile.ReadAt(bytes, int64(s.dimension*8*pos))
	if err != nil {
		panic(fmt.Errorf("Error reading vectors file"))
	}
	return BytesToArray(bytes)
}
