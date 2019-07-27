package store

import (
	"testing"
)

func TestSearch(t *testing.T) {
	//train vectors
	var v1 = Zeros("v1", 5)
	var v2 = Zeros("v2", 5)
	var v3 = Zeros("v3", 5)
	v3.Values()[3] = 1
	var v4 = Zeros("v4", 5)
	v4.Values()[3] = 5
	var v5 = Zeros("v5", 5)
	v5.Values()[3] = 2
	v5.Values()[4] = 1

	s := NewSimpleMapStore()
	if s.Set(v5.ID, v5) != nil || s.Set(v3.ID, v3) != nil || s.Set(v2.ID, v2) != nil || s.Set(v4.ID, v4) != nil || s.Set(v1.ID, v1) != nil {
		t.Error("error setting values")
	}

	neighbors, _ := s.KNN(v1, 5)
	if (*neighbors)[0].Target != v1 && (*neighbors)[0].Target != v2 {
		t.Error("wrong first value")
	}
	if (*neighbors)[1].Target != v1 && (*neighbors)[1].Target != v2 {
		t.Error("wrong second value")
	}
	if (*neighbors)[2].Target != v3 {
		t.Error("wrong third value")
	}
	if (*neighbors)[3].Target != v5 {
		t.Error("wrong fourth value")
	}
	if (*neighbors)[4].Target != v4 {
		t.Error("wrong fifth value")
	}

}
