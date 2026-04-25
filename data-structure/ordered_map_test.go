package main

import (
	"slices"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	keys := om.Keys()
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Errorf("expected keys [a b c], got %v", keys)
	}

	values := om.Values()
	if !slices.Equal(values, []int{1, 2, 3}) {
		t.Errorf("expected values [1 2 3], got %v", values)
	}

	value, exists := om.Get("b")
	if !exists || value != 2 {
		t.Errorf("expected to get value 2 for key 'b', got %v (exists: %v)", value, exists)
	}
}

func TestOrderedMapFrom(t *testing.T) {
	m := map[string]int{
		"c": 3,
		"a": 1,
		"b": 2,
	}
	om := OrderedMapFrom(m)

	keys := om.Keys()
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Errorf("expected keys [a b c], got %v", keys)
	}

	values := om.Values()
	if !slices.Equal(values, []int{1, 2, 3}) {
		t.Errorf("expected values [1 2 3], got %v", values)
	}
}

type pairs []struct {
	Key   string
	Value int
}

func TestOrderedMapFromPairs(t *testing.T) {
	pairs := pairs{
		{"c", 3},
		{"a", 1},
		{"b", 2},
	}
	om := OrderedMapFromPairs(pairs...)

	keys := om.Keys()
	if !slices.Equal(keys, []string{"c", "a", "b"}) {
		t.Errorf("expected keys [c a b], got %v", keys)
	}

	values := om.Values()
	if !slices.Equal(values, []int{3, 1, 2}) {
		t.Errorf("expected values [3 1 2], got %v", values)
	}
}
