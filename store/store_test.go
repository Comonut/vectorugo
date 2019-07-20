package store

import (
	"os"
	"testing"
)

func TestSum(t *testing.T) {
	var onesVector = Ones("", 32)
	if onesVector.Sum() != 32 {
		t.Errorf("Wrong sum")
	}
}

func set(s Store, id string, vector *Vector) error {
	return s.Set(id, vector)
}

func get(s Store, id string) (*Vector, error) {
	return s.Get(id)
}

func del(s Store, id string) error {
	return s.Delete(id)
}

func vectorEquals(this, other *Vector) bool {
	if this.ID != other.ID {
		return false
	}
	if len(this.Values) != len(other.Values) {
		return false
	}
	for i, e := range this.Values {
		if e != (*other).Values[i] {
			return false
		}
	}
	return true
}

func testStore(s Store, t *testing.T) {
	var ones = Ones("ones", 32)
	var zeros = Zeros("zeros", 32)
	var random = Random("rand", 32)

	if set(s, ones.ID, ones) != nil || set(s, zeros.ID, zeros) != nil || set(s, random.ID, random) != nil {
		t.Error("error setting values")
	}

	var val, err = get(s, ones.ID)
	if err != nil || !vectorEquals(val, ones) {
		t.Error("error getting values")
	}

	if del(s, ones.ID) != nil {
		t.Error("error deleting from store")
	}
	_, err = get(s, ones.ID)
	if err == nil {
		t.Error("can get value that should have been deleted")
	}
	val, err = get(s, random.ID)
	if err != nil || !vectorEquals(val, random) {
		t.Error("error getting values")
	}
	val, err = get(s, zeros.ID)
	if err != nil || !vectorEquals(val, zeros) {
		t.Error("error getting values")
	}
}

func TestSimpleMapStore(t *testing.T) {
	s := NewSimpleMapStore()
	testStore(&s, t)
}

func TestPersistantStore(t *testing.T) {
	os.Remove("index.test")
	os.Remove("vectors.test")
	s, _ := NewPersitantStore(uint32(32), "index.test", "vectors.test")
	testStore(s, t)
}
