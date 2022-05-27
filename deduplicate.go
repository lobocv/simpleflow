package simpleflow

// Deduplicator is an entity that keeps track of items it has seen before so that it can deduplicate values
type Deduplicator[T comparable] struct {
	seen map[T]struct{}
}

// NewDeduplicator returns a new Deduplicator which can be used to deduplicate slices values
func NewDeduplicator[T comparable]() *Deduplicator[T] {
	return &Deduplicator[T]{seen: make(map[T]struct{})}
}

// Add adds a item to the Deduplicator and returns true if it was a new value (ie not a duplicate)
func (dd *Deduplicator[T]) Add(v T) bool {
	_, exists := dd.seen[v]
	if exists {
		return false
	}
	dd.seen[v] = struct{}{}
	return true
}

// Seen returns true if the provided value has already been added to the Deduplicator
func (dd *Deduplicator[T]) Seen(v T) bool {
	_, exists := dd.seen[v]
	return exists
}

// Reset removes any memory of duplicate values seen by this Deduplicator{}
func (dd *Deduplicator[T]) Reset() {
	dd.seen = make(map[T]struct{})
}

// Deduplicate returns a newly allocated slice without duplicate values by comparing it against values previously
// seen by the Deduplicator{}
func (dd *Deduplicator[T]) Deduplicate(values []T) []T {
	var deduped []T
	for _, v := range values {
		if dd.Add(v) {
			deduped = append(deduped, v)
		}
	}
	return deduped
}

// DeduplicateIndices returns the indices of values in the provided slice which are duplicates
func (dd *Deduplicator[T]) DeduplicateIndices(values []T) []int {
	var indices []int
	for idx, v := range values {
		if !dd.Add(v) {
			indices = append(indices, idx)
		}
	}
	return indices
}

// Deduplicate returns a newly allocated slice without deplicate values
func Deduplicate[V comparable](values []V) []V {
	dd := NewDeduplicator[V]()
	return dd.Deduplicate(values)
}

// ObjectDeduplicator is a deduplicator that works on objects by creating an ID for each element. Objects
// with the same ID will be deduplicated.
type ObjectDeduplicator[T any] struct {
	dd   *Deduplicator[string]
	toId func(T) string
}

// NewObjectDeduplicator creates a ObjectDeduplicator that uses the provided function in order to create IDs for
// needing to be deduplicated.
func NewObjectDeduplicator[T any](toId func(T) string) *ObjectDeduplicator[T] {
	return &ObjectDeduplicator[T]{dd: NewDeduplicator[string](), toId: toId}
}

// Add adds a item to the ObjectDeduplicator and returns true if it was a new value (ie not a duplicate)
func (dd *ObjectDeduplicator[T]) Add(v T) bool {
	id := dd.toId(v)
	return dd.dd.Add(id)
}

// Seen returns true if the provided value has already been added to the ObjectDeduplicator
func (dd *ObjectDeduplicator[T]) Seen(v T) bool {
	id := dd.toId(v)
	return dd.dd.Seen(id)
}

// Reset removes any memory of duplicate values seen by this Deduplicator{}
func (dd *ObjectDeduplicator[T]) Reset() {
	dd.dd.Reset()
}

// Deduplicate returns a newly allocated slice without duplicate values by comparing it against values previously
// seen by the ObjectDuplicator{}
func (dd *ObjectDeduplicator[T]) Deduplicate(values []T) []T {
	var deduped []T
	for _, v := range values {
		if dd.Add(v) {
			deduped = append(deduped, v)
		}
	}
	return deduped
}

// DeduplicateIndices returns the indices of values in the provided slice which are duplicates
func (dd *ObjectDeduplicator[T]) DeduplicateIndices(values []T) []int {
	var indices []int
	for idx, v := range values {
		if !dd.Add(v) {
			indices = append(indices, idx)
		}
	}
	return indices
}
