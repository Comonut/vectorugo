package store

import (
	"encoding/binary"
	"fmt"
	"math"
)

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func ArrToBytes(arr []float64) []byte {
	bytes := make([]byte, 8*len(arr))
	var part []byte
	for i, f := range arr {
		part = Float64bytes(f)
		copy(bytes[i*8:i*8+8], part)
	}
	return bytes
}

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
