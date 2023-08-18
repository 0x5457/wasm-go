package wasm_go

type stack[T any] struct {
	inner []T
}

func (s *stack[T]) Len() int {
	return len(s.inner)
}

func (s *stack[T]) isEmpty() bool {
	return s.Len() == 0
}

func (s *stack[T]) Push(v T) {
	s.inner = append(s.inner, v)
}

func (s *stack[T]) Top() (*T, bool) {
	return s.Peek(0)
}

func (s *stack[T]) Peek(depth int) (*T, bool) {
	var v T
	if s.isEmpty() {
		return &v, false
	}
	return &s.inner[len(s.inner)-1-depth], true
}

func (s *stack[T]) Set(sp, idx int, v T) {
	s.inner[sp+idx] = v
}

func (s *stack[T]) Get(sp, idx int) (*T, bool) {
	if sp+idx >= s.Len() {
		return nil, false
	}
	return &s.inner[sp+idx], true
}

func (s *stack[T]) Pop() (T, bool) {
	var v T
	if s.isEmpty() {
		return v, false
	}
	idx := s.Len() - 1
	v = s.inner[idx]
	s.inner = s.inner[:idx]
	return v, true
}
