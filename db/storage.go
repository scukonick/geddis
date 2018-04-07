package db

import (
	"container/heap"
	"strings"
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

func (h *keyTimeHeap) deleteKey(key string) {
	var i int
	var val keyTime
	found := false
	for i, val = range *h {
		if val.key == key {
			found = true
			break
		}
	}

	if !found {
		return
	}

	heap.Remove(h, i)
}

// GeddisStore is Key Value storage.
// It can be used as embedded storage.
type GeddisStore struct {
	m      map[string]interface{}
	h      keyTimeHeap
	lock   sync.RWMutex
	stopCh chan interface{}
	wg     sync.WaitGroup
}

// NewGeddisStore returns newly initialized GeddisStore.
// It will allocate hash map enough to store 'size' elements.
// But it's not maximum, it's set only to prevent memory allocations.
func NewGeddisStore(size int) *GeddisStore {
	if size < 0 {
		size = 0
	}

	return &GeddisStore{
		m:      make(map[string]interface{}, size),
		h:      newKeyTimeHeap(size),
		lock:   sync.RWMutex{},
		stopCh: make(chan interface{}),
		wg:     sync.WaitGroup{},
	}
}

// cleanExpired removes expired elements from heap.
// It's not thread-safe so should be wrapped with mutex lock
func (s *GeddisStore) cleanExpired() {

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

// SetStr sets string value.
// If ttl <= 0, ttl is not set.
// If there is already an element with key = key, it
// would be overwritten
func (s *GeddisStore) SetStr(key, value string, ttl time.Duration) {
	s.set(key, value, ttl)
}

// SetArr sets string value.
// If ttl <= 0, ttl is not set.
// If there is already an element with key = key, it
// would be overwritten
func (s *GeddisStore) SetArr(key string, value []string, ttl time.Duration) {
	s.set(key, value, ttl)
}

// SetMap sets string value.
// If ttl <= 0, ttl is not set.
// If there is already an element with key = key, it
// would be overwritten
func (s *GeddisStore) SetMap(key string, value map[string]string, ttl time.Duration) {
	s.set(key, value, ttl)
}

func (s *GeddisStore) set(key string, value interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.m[key]; exists {
		// cleaning heap from old TTL
		s.h.deleteKey(key)
	}

	s.m[key] = value
	if ttl != 0 {
		expireAt := time.Now().Add(ttl)
		heap.Push(&s.h, newKeyTime(key, expireAt))
	}
}

// GetStr returns string stored by key 'key'.
// If an element with this key does not exist, it returns ErrNotFound error.
// If by this key there is value of another type,
// it returns ErrInvalidType.
func (s *GeddisStore) GetStr(key string) (string, error) {
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

// GetArr returns a slice of strings stored by key 'key'.
// If an element with this key does not exist, it returns ErrNotFound error.
// If by this key there is value of another type,
// it returns ErrInvalidType.
func (s *GeddisStore) GetArr(key string) ([]string, error) {
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

// GetMap returns a map stored by key 'key'.
// If an element with this key does not exist, it returns ErrNotFound error.
// If by this key there is value of another type,
// it returns ErrInvalidType.
func (s *GeddisStore) GetMap(key string) (map[string]string, error) {
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

func (s *GeddisStore) del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.m, key)

	s.h.deleteKey(key)
}

// Keys returns slice of keys which have prefix 'prefix'.
// If prefix is an empty string, it returns all the keys
func (s *GeddisStore) Keys(prefix string) []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp := make([]string, 0, len(s.m))
	for k := range s.m {
		if strings.HasPrefix(k, prefix) {
			resp = append(resp, k)
		}
	}

	return resp
}

// GetByIndex get element by index i from slice
// by key 'key'. It returns ErrNotFound if such slice
// does not exist, ErrInvalidType if element by this key
// is not a slice and ErrInvalidIndex if element does not contain
// sub-element defined by this index. It also returns ErrInvalidIndex
// in case of i < 0.
func (s *GeddisStore) GetByIndex(key string, i int) (string, error) {
	if i < 0 {
		return "", ErrInvalidIndex
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}

	elemSlice, ok := elem.([]string)
	if !ok {
		return "", ErrInvalidType
	}

	if len(elemSlice) < i+1 {
		return "", ErrInvalidIndex
	}

	return elemSlice[i], nil
}

// GetByKey get element by key subKey from map
// by key 'key'. It returns ErrNotFound if such map
// does not exist, ErrInvalidType if element by this key
// is not a map and ErrInvalidIndex if element does not contain
// sub-element defined by this key.
func (s *GeddisStore) GetByKey(key, subKey string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.cleanExpired()

	elem, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}

	elemSlice, ok := elem.(map[string]string)
	if !ok {
		return "", ErrInvalidType
	}

	resp, ok := elemSlice[subKey]
	if !ok {
		return "", ErrInvalidIndex
	}

	return resp, nil
}

// Run starts inner processes of strings store
// i.e. cleaning expired elements and storing data to disk.
// It starts processes in separate goroutines and exits after it.
func (s *GeddisStore) Run() {
	go s.runCleaner()
}

// runCleaner cleans expired elements each 5 seconds
// in order to free memory when no one is calling 'get' method
func (s *GeddisStore) runCleaner() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Minute)

		for {
			select {
			case <-ticker.C:
				s.lock.Lock()
				s.cleanExpired()
				s.lock.Unlock()
			case <-s.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops all processes of store and waits for them to exit.
func (s *GeddisStore) Stop() {
	s.stopCh <- true

	s.wg.Wait()
}
