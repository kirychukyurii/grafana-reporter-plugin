package cron

import "sync"

// Schedulers map scheduler per organization
type Schedulers struct {
	mx sync.RWMutex
	m  map[int]SchedulerManager
}

func (s *Schedulers) Load(key int) (SchedulerManager, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.m[key]

	return val, ok
}

func (s *Schedulers) Store(key int, value SchedulerManager) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.m[key] = value
}
