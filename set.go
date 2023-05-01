package maps

// Set is a map[T]struct{}
type Set[T comparable] map[T]struct{}

// NewSet creates an empty Set[T]
func NewSet[T comparable]() Set[T] { return Set[T](make(map[T]struct{})) }

// Has checks value is in Set
func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}

// Add stores a value
func (s Set[T]) Add(t T) { s[t] = struct{}{} }

// Remove deletes a value
func (s Set[T]) Remove(t T) { delete(s, t) }

// Slice returns this Set[T] as []T
func (s Set[T]) Slice() []T {
	i, slice := 0, make([]T, len(s))
	for v := range s {
		slice[i] = v
		i++
	}
	return slice
}

// Each calls a function once for every value
func (s Set[T]) Each(f func(v T)) {
	for v := range s {
		f(v)
	}
}

// Delete deletes items
func (s Set[T]) Delete(items ...T) {
	for _, t := range items {
		delete(s, t)
	}
}
