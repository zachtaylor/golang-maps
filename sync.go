package maps

import "sync"

// Sync is a generic RWMutex map
type Sync[K comparable, V any] struct {
	rw   sync.RWMutex
	data map[K]V
}

// NewSync creates an empty *Sync[K, V]
func NewSync[K comparable, V any]() *Sync[K, V] {
	return &Sync[K, V]{
		data: make(map[K]V),
	}
}

// Keys returns the keys
func (s *Sync[K, V]) Keys() []K {
	s.rw.RLock()
	i, keys := 0, make([]K, len(s.data))
	for k := range s.data {
		keys[i] = k
		i++
	}
	s.rw.RUnlock()
	return keys
}

// Values returns the values
func (s *Sync[K, V]) Values() []V {
	s.rw.RLock()
	i, values := 0, make([]V, len(s.data))
	for _, v := range s.data {
		values[i] = v
		i++
	}
	s.rw.RUnlock()
	return values
}

// Size returns the number of items
func (s *Sync[K, V]) Size() int { return len(s.data) }

// Get returns the value for a key
func (s *Sync[K, V]) Get(key K) V { return s.data[key] }

// Set changes the value for a key
func (s *Sync[K, V]) Set(key K, val V) {
	s.rw.Lock()
	s.set(key, val)
	s.rw.Unlock()
}

func (s *Sync[K, V]) set(key K, val V) {
	s.data[key] = val
}

// Clone returns a shallow clone of a map
func (s *Sync[K, V]) Clone() *Sync[K, V] {
	if s == nil {
		return nil
	}
	ss := NewSync[K, V]()
	s.rw.RLock()
	for k, v := range s.data {
		ss.data[k] = v
	}
	s.rw.RUnlock()
	return ss
}

// Each calls a function, once for every value, inside the mutex lock state
func (s *Sync[K, V]) Each(f func(K, V)) {
	s.rw.RLock()
	for k, v := range s.data {
		f(k, v)
	}
	s.rw.RUnlock()
}

// Filter uses a test func to filter the map
func (s *Sync[K, V]) Filter(f func(K, V) bool) map[K]V {
	filtered := make(map[K]V)
	s.rw.RLock()
	for k, v := range s.data {
		if f(k, v) {
			filtered[k] = v
		}
	}
	s.rw.RUnlock()
	return filtered
}

// Find uses a test func to find the first passing value
func (s *Sync[K, V]) Find(f func(K, V) bool) (_ K, _ V) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for k, v := range s.data {
		if f(k, v) {
			return k, v
		}
	}
	return
}

// ReduceSync returns an accumulation of a *Sync using an accumulation func
func ReduceSync[K comparable, V any, A any](s *Sync[K, V], a A, f func(A, K, V) A) A {
	if s == nil {
		return a
	}
	s.rw.RLock()
	for k, v := range s.data {
		a = f(a, k, v)
	}
	s.rw.RUnlock()
	return a
}

// Delete deletes keys
func (s *Sync[K, V]) Delete(keys ...K) {
	s.rw.Lock()
	for _, key := range keys {
		delete(s.data, key)
	}
	s.rw.Unlock()
}

// DeleteFunc deletes where del returns true
func (s *Sync[K, V]) DeleteFunc(del func(K, V) bool) {
	s.rw.Lock()
	for k, v := range s.data {
		if del(k, v) {
			delete(s.data, k)
		}
	}
	s.rw.Unlock()
}

// Lock calls a function inside the RWMutex write lock state
func (s *Sync[K, V]) Lock(f func(set func(K, V))) {
	s.rw.Lock()
	f(s.set)
	s.rw.Unlock()
}

// RLock calls a function inside the RWMutex read lock state
func (s *Sync[K, V]) RLock(f func()) {
	s.rw.RLock()
	f()
	s.rw.RUnlock()
}
