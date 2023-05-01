package maps

// Keys returns the keys of a map
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// Values returns the values of a map
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// Clone returns a shallow clone of a map
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

func Each[M ~map[K]V, K comparable, V any](m M, f func(K, V)) {
	if m == nil {
		return
	}
	for k, v := range m {
		f(k, v)
	}
}

// Filter returns a shallow clone of a map containing each entry where test returns true
func Filter[M ~map[K]V, K comparable, V any](m M, test func(K, V) bool) M {
	if m == nil {
		return nil
	}
	r := make(M)
	for k, v := range m {
		if test(k, v) {
			r[k] = v
		}
	}
	return r
}

// Find returns the entry key and value for the first entry where test returns true
func Find[M ~map[K]V, K comparable, V any](m M, test func(K, V) bool) (_k K, _v V) {
	if m == nil {
		return
	}
	for k, v := range m {
		if test(k, v) {
			_k, _v = k, v
			break
		}
	}
	return
}

// Reduce returns an accumulation of a map using an accumulation func
func Reduce[M ~map[K]V, K comparable, V any, A any](m M, a A, f func(A, K, V) A) A {
	if m == nil {
		return a
	}
	for k, v := range m {
		a = f(a, k, v)
	}
	return a
}

// Copy writes all key/value pairs in src to dst
func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

// DeleteFunc deletes from a map where del returns true
func DeleteFunc[M ~map[K]V, K comparable, V any](m M, del func(K, V) bool) {
	for k, v := range m {
		if del(k, v) {
			delete(m, k)
		}
	}
}
