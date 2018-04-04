package db

import (
	"container/heap"
	"sync"
	"time"
)

type stringTime struct {
	str      string
	expireAt time.Time
}

func newStringKey(key string, expireAt time.Time) stringTime {
	return stringTime{
		str:      key,
		expireAt: expireAt,
	}
}

type stringKeyHeap []stringTime

func newStringKeyHeap(size int) stringKeyHeap {
	return make([]stringTime, 0, size)
}

func (h stringKeyHeap) Len() int {
	return len(h)
}
func (h stringKeyHeap) Less(i, j int) bool {
	return h[i].expireAt.Before(h[j].expireAt)
}
func (h stringKeyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *stringKeyHeap) Push(x interface{}) {
	*h = append(*h, x.(stringTime))
}

func (h *stringKeyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// stringsStore represents string-value storage
// for storing strings
type stringsStore struct {
	m      map[string]string
	h      stringKeyHeap
	lock   sync.RWMutex
	stopCh chan interface{}
	wg     sync.WaitGroup
}

func newStringsStore(size int) *stringsStore {
	if size < 0 {
		size = 0
	}

	return &stringsStore{
		m:      make(map[string]string, size),
		h:      newStringKeyHeap(size),
		lock:   sync.RWMutex{},
		stopCh: make(chan interface{}),
		wg:     sync.WaitGroup{},
	}
}

// cleanExpired removes expired elements from heap.
// It's not thread safe so should be wrapped with mutex lock
func (s *stringsStore) cleanExpired() {

	now := time.Now()

	for {
		if len(s.h) == 0 {
			break
		}
		if s.h[0].expireAt.Before(now) {
			k := heap.Pop(&s.h)
			delete(s.m, k.(stringTime).str)
		}
	}
}

func (s *stringsStore) set(key, value string, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.m[key] = value
	if ttl != 0 {
		expireAt := time.Now().Add(ttl)
		s.h.Push(newStringKey(key, expireAt))
	}
}

func (s *stringsStore) get(key string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}

	return elem, nil
}

func (s *stringsStore) del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.m, key)

	// instead of ranging over all the heap
	// we can store ttl also in the map
	// but it will increase memory usage.
	// so here we prefer some more CPU usage
	// in order to save memory
	for i, v := range s.h {
		if v.str == key {
			heap.Remove(&s.h, i)
			break
		}
	}
}

// runCleaner cleans expired elements each 5 seconds
// to free memory when no one is calling 'get' method
func (s *stringsStore) runCleaner() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for {
			tick := time.NewTimer(1 * time.Minute)
			select {
			case <-tick.C:
				s.lock.Lock()
				s.cleanExpired()
				s.lock.Unlock()
			case <-s.stopCh:
				tick.Stop()
				return
			}
		}
	}()
}
