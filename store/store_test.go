package store

import "testing"

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
func testStore(s Store, t *testing.T) {
	var ones = Ones("ones", 32)
	var zeros = Zeros("zeros", 32)


	if set(s, ones.ID, ones) != nil || set(s, zeros.ID, zeros) != nil {
		t.Error("error setting values")
	}

	var val, err = get(s, ones.ID)
	if err != nil || val != ones {
		t.Error("error getting values")
	}

	if del(s, ones.ID) != nil {
		t.Error("error deleting from store")
	}
	_, err = get(s, ones.ID)
	if err == nil {
		t.Error("can get value that should have been deleted")
	}
}

func TestSimpleMapStore(t *testing.T) {
	var s = NewSimpleMapStore()
	testStore(&s, t)
}
