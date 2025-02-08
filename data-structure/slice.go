package main

type Slice[T any] struct {
	arr         []T
	cap         int
	len         int
	extendRatio int // 扩容倍数
}

func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{
		arr:         make([]T, 10),
		cap:         10, // 列表容量
		len:         0,  // 列表长度 当前元素的数量
		extendRatio: 2,  // 每次列表扩容的倍数
	}
}

func (s *Slice[T]) Len() int {
	return s.len
}

func (s *Slice[T]) Cap() int {
	return s.cap
}

func (s *Slice[T]) Get(index int) T {
	if index < 0 || index >= s.len {
		panic("index out of range")
	}
	return s.arr[index]
}

func (s *Slice[T]) Set(index int, val T) {
	if index < 0 || index >= s.len {
		panic("index out of range")
	}
	s.arr[index] = val
}

func (s *Slice[T]) Append(val T) {
	if s.len == s.cap {
		s.extend()
	}
	s.arr[s.len] = val
	s.len++
}

func (s *Slice[T]) Insert(index int, val T) {
	if index < 0 || index > s.len {
		panic("index out of range")
	}
	if s.len == s.cap {
		s.extend()
	}
	for i := s.len; i > index; i-- {
		s.arr[i] = s.arr[i-1]
	}
	s.arr[index] = val
	s.len++
}

func (s *Slice[T]) Delete(index int) T {
	if index < 0 || index >= s.len {
		panic("index out of range")
	}
	e := s.arr[index]
	for i := index; i < s.len-1; i++ {
		s.arr[i] = s.arr[i+1]
	}
	s.len--
	return e // 返回被删除的元素
}

func (s *Slice[T]) extend() {
	if s.cap == 0 {
		s.cap = 1
	} else {
		s.cap *= s.extendRatio
	}
	newArr := make([]T, s.cap)
	copy(newArr, s.arr)
	s.arr = newArr
}
