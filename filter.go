package cuckoofilter

import "fmt"

type Filter[T Hash] struct {
	hashtable [][]int
	pos       []int
	ver       int
	length    int
}

func NewFilter[T Hash](ver, length int) *Filter[T] {
	hashtable := make([][]int, ver)
	for i := 0; i < ver; i++ {
		hashtable[i] = make([]int, length)
	}

	pos := make([]int, ver)

	return &Filter[T]{
		pos:       pos,
		hashtable: hashtable,
		ver:       ver,
		length:    length,
	}
}

func (s *Filter[T]) Add(item T) error {
	return s.add(item.Sum(), 0, 0)
}

func (s *Filter[T]) add(key, tableID, cnt int) error {
	if cnt >= s.length {
		return fmt.Errorf("%v unpositioned", key)
	}

	ok := s.query(key)
	if ok {
		return nil
	}

	if s.hashtable[tableID][s.pos[tableID]] != 0 {
		dis := s.hashtable[tableID][s.pos[tableID]]
		s.hashtable[tableID][s.pos[tableID]] = key
		return s.add(dis, (tableID+1)%s.ver, cnt+1)
	} else {
		s.hashtable[tableID][s.pos[tableID]] = key
	}

	return nil
}

func (s *Filter[T]) Contains(item T) bool {
	key := item.Sum()
	return s.query(key)
}

func (s *Filter[T]) Remove(item T) {
	key := item.Sum()

	for i := 0; i < s.ver; i++ {
		s.pos[i] = s.position(i+1, key)
		if s.hashtable[i][s.pos[i]] == key {
			s.hashtable[i][s.pos[i]] = 0
		}
	}
}

func (s *Filter[T]) query(key int) bool {
	for i := 0; i < s.ver; i++ {
		s.pos[i] = s.position(i+1, key)
		if s.hashtable[i][s.pos[i]] == key {
			return true
		}
	}
	return false
}

func (s *Filter[T]) position(i int, key int) int {
	switch i {
	case 1:
		return key % s.length
	case 2:
		return (key / s.length) % s.length
	}
	return 0
}
