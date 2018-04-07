package db

import (
	"container/heap"
	"sync"
	"time"
)

type keyTime struct {
	key      string
	expireAt time.Time
}

func newKeyTime(key string, expireAt time.Time) keyTime {
	return keyTime{
		key:      key,
		expireAt: expireAt,
	}
}

type keyTimeHeap []keyTime

func newKeyTimeHeap(size int) keyTimeHeap {
	return make([]keyTime, 0, size)
}

func (h keyTimeHeap) Len() int {
	return len(h)
}
func (h keyTimeHeap) Less(i, j int) bool {
	return h[i].expireAt.Before(h[j].expireAt)
}
func (h keyTimeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *keyTimeHeap) Push(x interface{}) {
	*h = append(*h, x.(keyTime))
}

func (h *keyTimeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// stringsStore represents string-value storage
// for storing strings
type stringsStore struct {
	m      map[string]interface{}
	h      keyTimeHeap
	lock   sync.RWMutex
	stopCh chan interface{}
	wg     sync.WaitGroup
}

func newStringsStore(size int) *stringsStore {
	if size < 0 {
		size = 0
	}

	return &stringsStore{
		m:      make(map[string]interface{}, size),
		h:      newKeyTimeHeap(size),
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
			delete(s.m, k.(keyTime).key)
		}
	}
}

func (s *stringsStore) set(key string, value interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.m[key] = value
	if ttl != 0 {
		expireAt := time.Now().Add(ttl)
		s.h.Push(newKeyTime(key, expireAt))
	}
}

// we could not use get for everything here as we need to copy values
// for slices and maps to prevent racing
// so we have to have different get's for each type
func (s *stringsStore) getStr(key string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}

	resp, ok := elem.(string)
	if !ok {
		return "", ErrInvalidType
	}

	return resp, nil
}

func (s *stringsStore) getArr(key string) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return nil, ErrNotFound
	}

	elemSlice, ok := elem.([]string)
	if !ok {
		return nil, ErrInvalidType
	}

	resp := make([]string, len(elemSlice))
	copy(resp, elemSlice)

	return resp, nil
}

func (s *stringsStore) getMap(key string) (map[string]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return nil, ErrNotFound
	}

	elemMap, ok := elem.(map[string]string)
	if !ok {
		return nil, ErrInvalidType
	}

	resp := make(map[string]string, len(elemMap))
	for i, val := range elemMap {
		resp[i] = val
	}

	return resp, nil
}

func (s *stringsStore) del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.m, key)

	for i, v := range s.h {
		if v.key == key {
			heap.Remove(&s.h, i)
			break
		}
	}
}

// runCleaner cleans expired elements each 5 seconds
// in order to free memory when no one is calling 'get' method
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

func (s *stringsStore) stop() {
	s.stopCh <- true

	s.wg.Wait()
}
