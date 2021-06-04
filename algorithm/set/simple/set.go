package simple

// 最简单的set实现

type set struct {
	m map[string]struct{}
}

func NewSet() *set {
	return &set{m: make(map[string]struct{})}
}

func (s *set) Has(key string) bool {
	_, ok := s.m[key]
	return ok
}

func (s *set) Add(key string) {
	s.m[key] = struct{}{}
}

func (s *set) Delete(key string) {
	delete(s.m, key)
}

func (s *set) Exist(k string) bool {
	_, ok := s.m[k]
	return ok
}
