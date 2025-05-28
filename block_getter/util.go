package block_getter

import "sync"

type SafeVar[T any] struct {
	Value T
	Lock  sync.RWMutex
}

func (s *SafeVar[T]) Set(value T) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Value = value
}

func (s *SafeVar[T]) Get() T {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.Value
}
