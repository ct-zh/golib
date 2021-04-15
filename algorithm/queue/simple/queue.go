package simple

import "sync"

type Queue interface {
	LPush(v interface{})
	LPop() interface{}
	RPush(v interface{})
	RPop() interface{}
	IsEmpty() bool
	Len() int
}

type simpleQueue []interface{}

func New() Queue {
	return &simpleQueue{}
}

func (q *simpleQueue) RPush(v interface{}) {
	*q = append(*q, v)
}

func (q *simpleQueue) RPop() interface{} {
	tail := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return tail
}

func (q *simpleQueue) LPush(v interface{}) {
	head := []interface{}{v}
	*q = append(head, *q...)
}

func (q *simpleQueue) LPop() interface{} {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

func (q *simpleQueue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *simpleQueue) Len() int {
	return len(*q)
}

// 附带锁的queue
type syncQueue struct {
	mu sync.Mutex
	q  []interface{}
}

func NewSyncQueue() Queue {
	return &syncQueue{}
}

func (s *syncQueue) LPush(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.q = append([]interface{}{v}, s.q...)
}

func (s *syncQueue) LPop() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	head := s.q[0]
	s.q = s.q[1:]
	return head
}

func (s *syncQueue) RPush(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.q = append(s.q, v)
}

func (s *syncQueue) RPop() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	tail := s.q[len(s.q)-1]
	s.q = s.q[0 : len(s.q)-1]
	return tail
}

func (s *syncQueue) IsEmpty() bool {
	return len(s.q) == 0
}

func (s *syncQueue) Len() int {
	return len(s.q)
}
