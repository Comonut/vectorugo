package store

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/sirupsen/logrus"
)

const mask = "@@"

type Index struct {
	maxlen   int
	size     int64
	branches []*branch
	db       *leveldb.DB
}

type branch struct {
	pos   Vector
	leafs []Vector
}

func NewIndex(db *leveldb.DB, store Store) *Index {
	if db != nil {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), db: db}
	} else {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), db: nil}
	}
}
func LoadIndex(db *leveldb.DB, s *PersistantStore) *Index {
	branchesArr := make([]*branch, 0)
	branchesMap := make(map[uint32]int)
	var leaf Vector

	var err error

	iter := db.NewIterator(util.BytesPrefix([]byte("@@")), nil)
	for iter.Next() {
		// Use key/value.
		leaf := 
	}
	iter.Release()

	for i := 0; i < len(inversePosIndex); i++ {
		_, err = file.ReadAt(leafPosBytes, int64(i*8))
		if err != nil {
			logrus.Error("error loading index leaf")
			panic(err)
		}
		_, err = file.ReadAt(branchPosBytes, int64(i*8+4))
		if err != nil {
			logrus.Error("error loading index branch")
			panic(err)
		}

		leafPos = binary.LittleEndian.Uint32(leafPosBytes)
		branchPos = binary.LittleEndian.Uint32(branchPosBytes)

		leaf = &PersistantVector{inversePosIndex[leafPos], leafPos, s}

		branchID, ok := branchesMap[branchPos]

		if ok {
			branchesArr[branchID].leafs = append(branchesArr[branchID].leafs, leaf)
		} else {
			branchesMap[branchPos] = len(branchesArr)
			branchesArr = append(branchesArr, &branch{pos: &MemoryVector{ID: inversePosIndex[branchPos], Array: s.ReadAtPos(branchPos)}, id: branchPos, leafs: []Vector{leaf}})
		}
	}

	return &Index{maxlen: int(2 * (math.Sqrt(float64(len(inversePosIndex))))), size: int64(len(inversePosIndex)), branches: branchesArr, file: file}

}

func (index *Index) writeLeafToFile(leaf Vector, branch *branch) {
	if index.db == nil {
		return
	}
	leafBytes := append([]byte(mask), []byte(leaf.Name())...)
	branchBytes := append([]byte(mask), []byte(branch.pos.Name())...)

	err := index.db.Put(leafBytes, branchBytes, nil)
	if err != nil {
		logrus.Error("error writing index changes to file")
		panic(err)
	}
}

func (index *Index) transfer(old, new *branch) {
	newLeafs := make([]Vector, 0)
	updatedOldLeafs := make([]Vector, 0)

	for _, l := range old.leafs {
		if l == old.pos {
			updatedOldLeafs = append(updatedOldLeafs, l)
		} else if EuclideanDistance(l, old.pos) < EuclideanDistance(l, new.pos) {
			updatedOldLeafs = append(updatedOldLeafs, l)
		} else {
			newLeafs = append(newLeafs, l)
			index.writeLeafToFile(l, new)
		}
	}

	old.leafs = updatedOldLeafs
	new.leafs = append(new.leafs, newLeafs...)
}

func (index *Index) AddVector(v Vector) {
	if len(index.branches) == 0 {
		new := &branch{pos: &MemoryVector{ID: v.Name(), Array: *v.Values()}, leafs: []Vector{v}}
		index.branches = append(index.branches, new)
		if index.db != nil {
			index.writeLeafToFile(v, new)
		}
		index.size++
		return
	}

	var closest *branch
	var closestDist = math.MaxFloat64

	for _, b := range index.branches {
		if EuclideanDistance(b.pos, v) < closestDist {
			closest = b
			closestDist = EuclideanDistance(b.pos, v)
		}
	}

	if len(closest.leafs) == index.maxlen {
		new := &branch{pos: &MemoryVector{ID: v.Name(), Array: *v.Values()}, leafs: []Vector{v}}
		index.branches = append(index.branches, new)
		if index.db != nil {
			index.writeLeafToFile(v, new)
		}
		index.transfer(closest, index.branches[len(index.branches)-1])
	} else {
		closest.leafs = append(closest.leafs, v)
		if index.db != nil {
			index.writeLeafToFile(v, closest)
		}
	}
	index.size++
}

type BranchDistance struct {
	Target   *branch
	Distance float64 //Distance
}

func (index *Index) IndexKNN(k int, v Vector) (*[]Distance, error) {
	if int64(k) > index.size {
		return nil, fmt.Errorf("can't retrieve %d results - only %d are indexed", k, index.size)
	}

	sortedBranches := make([]BranchDistance, len(index.branches))

	//loop through map, calculate distance for each vector, append result in return array
	for i, b := range index.branches {
		sortedBranches[i] = BranchDistance{Target: b, Distance: EuclideanDistance(b.pos, v)}
	}
	//sort array by distance in incrementing order
	sort.Slice(sortedBranches, func(i, j int) bool {
		return sortedBranches[i].Distance < sortedBranches[j].Distance
	})

	potentials := make([]Vector, 0)

	for i := 0; len(potentials) < k; i++ {
		potentials = append(potentials, sortedBranches[i].Target.leafs...)
	}

	results := make([]Distance, len(potentials))

	for i, p := range potentials {
		results[i] = Distance{Target: p, Distance: EuclideanDistance(p, v)}
	}
	//sort array by distance in incrementing order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
	results = results[:k]
	return &results, nil

}
