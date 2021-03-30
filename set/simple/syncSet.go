package simple

import "sync"

// 支持并发的set

type SyncSet struct {
	m sync.Map
}

func (s *SyncSet) Has(key string) bool {
	_, ok := s.m.Load(key)
	return ok
}

func (s *SyncSet) Add(key string) {
	s.m.Store(key, struct{}{})
}

func (s *SyncSet) Delete(key string) {
	s.m.Delete(key)
}
