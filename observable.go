package maps

// Observer is an interface for observing a generic type
type Observer[K comparable, V any] interface {
	Observe(id K, new, old V)
}

// ObserverFunc is a func type that implements Observer
type ObserverFunc[K comparable, V any] func(id K, new, old V)

func (f ObserverFunc[K, V]) Observe(id K, new, old V) { f(id, new, old) }

// Observable is a generic observable map
type Observable[K comparable, V any] struct {
	sync Sync[K, V]
	obs  []Observer[K, V]
}

// NewObservable creates an empty *Observable[K, V]
func NewObservable[K comparable, V any]() *Observable[K, V] {
	return &Observable[K, V]{
		sync: Sync[K, V]{data: make(map[K]V)},
		obs:  make([]Observer[K, V], 0),
	}
}

func (o *Observable[K, V]) callback(key K, new, old V) {
	for _, o := range o.obs {
		o.Observe(key, new, old)
	}
}
func (o *Observable[K, V]) set(key K, val V) {
	o.callback(key, val, o.sync.data[key])
	o.sync.data[key] = val
}

// Keys returns the keys
func (o *Observable[K, V]) Keys() []K {
	return o.sync.Keys()
}

// Values returns the values
func (o *Observable[K, V]) Values() []V {
	return o.sync.Values()
}

// Size returns the number of items
func (o *Observable[K, V]) Size() int {
	return o.sync.Size()
}

// Get returns the value for a key
func (o *Observable[K, V]) Get(key K) V {
	return o.sync.Get(key)
}

// Set changes the value for a key
func (o *Observable[K, V]) Set(key K, val V) {
	o.sync.rw.Lock()
	o.set(key, val)
	o.sync.rw.Unlock()
}

// Observe adds an observer
func (o *Observable[K, V]) Observe(f Observer[K, V]) {
	o.obs = append(o.obs, f)
}

// Each calls a function, once for every value, inside the mutex lock state
func (o *Observable[K, V]) Each(f func(K, V)) {
	o.sync.Each(f)
}

// Filter uses a test func to filter the map
func (o *Observable[K, V]) Filter(f func(K, V) bool) map[K]V {
	return o.sync.Filter(f)
}

// Find uses a test func to find the first passing value
func (o *Observable[K, V]) Find(f func(K, V) bool) (key K, val V) {
	return o.sync.Find(f)
}

// ReduceSync returns an accumulation of an *Observable using an accumulation func
func ReduceObservable[K comparable, V any, A any](o *Observable[K, V], a A, f func(A, K, V) A) A {
	return ReduceSync(&o.sync, a, f)
}

// Delete deletes keys
func (o *Observable[K, V]) Delete(keys ...K) {
	var zero V
	o.sync.rw.Lock()
	for _, key := range keys {
		o.callback(key, zero, o.sync.data[key])
		delete(o.sync.data, key)
	}
	o.sync.rw.Unlock()
}

// DeleteFunc deletes where del returns true
func (o *Observable[K, V]) DeleteFunc(del func(K, V) bool) {
	var zero V
	o.sync.rw.Lock()
	for k, v := range o.sync.data {
		if del(k, v) {
			o.callback(k, zero, v)
			delete(o.sync.data, k)
		}
	}
	o.sync.rw.Unlock()
}

// Lock calls a function inside the RWMutex write lock state
func (o *Observable[K, V]) Lock(f func(set func(K, V))) {
	o.sync.rw.Lock()
	f(o.set)
	o.sync.rw.Unlock()
}

// RLock calls a function inside the RWMutex read lock state
func (o *Observable[K, V]) RLock(f func()) { o.sync.RLock(f) }
