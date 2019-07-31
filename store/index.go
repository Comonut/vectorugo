package store

import (
	"encoding/binary"
	"math"
	"os"
	"sort"
)

type Index struct {
	maxlen   int
	size     int64
	branches []*branch
	file     *os.File
}

type branch struct {
	pos   Vector
	leafs []Vector
}

func NewIndex(file *os.File) *Index {
	if file != nil {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), file: file}
	} else {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), file: nil}
	}
}
func LoadIndex(file *os.File) *Index {
	if file != nil {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), file: file}
	} else {
		return &Index{maxlen: 2, size: 0, branches: make([]*branch, 0), file: nil}
	}
}

func (index *Index) writeLeafToFile(leaf, branch Vector) {

	leafPv := leaf.(*PersistantVector)
	branchPv := branch.(*PersistantVector)

	leafBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(leafBytes, leafPv.pos)
	branchBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(branchBytes, branchPv.pos)

	allBytes := append(leafBytes, branchBytes...)

	index.file.WriteAt(allBytes, index.size)
}

func (index *Index) overrideLeafToFile(leaf, branch Vector) {

	leafPv := leaf.(*PersistantVector)
	branchPv := branch.(*PersistantVector)

	leafBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(leafBytes, leafPv.pos)
	branchBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(branchBytes, branchPv.pos)

	allBytes := append(leafBytes, branchBytes...)

	index.file.WriteAt(allBytes, int64(leafPv.pos))
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
			index.overrideLeafToFile(l, new.pos)
		}
	}

	old.leafs = updatedOldLeafs
	new.leafs = append(new.leafs, newLeafs...)
}

func (index *Index) AddVector(v Vector) {
	if len(index.branches) == 0 {
		index.branches = append(index.branches, &branch{pos: &MemoryVector{ID: v.Name(), Array: *v.Values()}, leafs: []Vector{v}})
		if index.file != nil {
			index.writeLeafToFile(v, v)
		}
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
		index.branches = append(index.branches, &branch{pos: &MemoryVector{ID: v.Name(), Array: *v.Values()}, leafs: []Vector{v}})
		if index.file != nil {
			index.writeLeafToFile(v, v)
		}
		index.transfer(closest, index.branches[len(index.branches)-1])
	} else {
		closest.leafs = append(closest.leafs, v)
		if index.file != nil {
			index.writeLeafToFile(v, closest.pos)
		}
	}
	index.size += 1
}

type BranchDistance struct {
	Target   *branch
	Distance float64 //Distance
}

func (index *Index) IndexKNN(k int, v Vector) *[]Distance {
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
	return &results

}
