package main

import "iter"

// Set holds a set of elements.
type Set[E comparable] struct {
	m map[E]struct{}
}

// New returns a new [Set]
func NewSet[E comparable]() *Set[E] {
	return &Set[E]{
		m: make(map[E]struct{}),
	}
}

// Add adds an element to a set
func (s *Set[E]) Add(v E) {
	s.m[v] = struct{}{}
}

// Contains reports whether an element is in a set.
func (s *Set[E]) Contains(v E) bool {
	_, ok := s.m[v]
	return ok
}

// Union returns the union of two sets
func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
	r := NewSet[E]()
	for v := range s1.m {
		r.Add(v)
	}
	for v := range s2.m {
		r.Add(v)
	}
	return r
}

// All is an iterator over the elements of s
func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}
