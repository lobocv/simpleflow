package simpleflow

// Counter is an entity that keeps track of the number items it encounters
type Counter[T comparable] struct {
	counts map[T]int
}

// NewCounter returns a new Counter which can be used to deduplicate slices values
func NewCounter[T comparable]() *Counter[T] {
	return &Counter[T]{counts: make(map[T]int)}
}

// Add adds a item to the Counter and returns the current number of occurrences
func (c *Counter[T]) Add(v T) int {
	c.counts[v]++
	return c.counts[v]
}

// Count returns the current number of occurrences for the given value
func (c *Counter[T]) Count(v T) int {
	return c.counts[v]
}

// Reset clears the values in the Counter{}
func (c *Counter[T]) Reset() {
	c.counts = make(map[T]int)
}

// AddMany adds all the values in the provided slice to the counter
func (c *Counter[T]) AddMany(values []T) {
	for _, v := range values {
		c.Add(v)
	}
}

// ObjectCounter is a counter that works on objects by creating an ID for each element. Objects
// with the same ID will be counted in the same bucket.
type ObjectCounter[T any] struct {
	c    *Counter[string]
	toId func(T) string
}

// NewObjectCounter creates a ObjectCounter that uses the provided function in order to create IDs for
// needing to be counted.
func NewObjectCounter[T any](toId func(T) string) *ObjectCounter[T] {
	return &ObjectCounter[T]{c: NewCounter[string](), toId: toId}
}

// Add adds an object to the Counter and returns the current number of occurrences
func (c *ObjectCounter[T]) Add(v T) int {
	return c.c.Add(c.toId(v))
}

// Count returns the current number of occurrences for the given object
func (c *ObjectCounter[T]) Count(v T) int {
	return c.c.Count(c.toId(v))
}

// Reset clears the values in the ObjectCounter{}
func (c *ObjectCounter[T]) Reset() {
	c.c.Reset()
}

// AddMany adds all the values in the provided slice to the counter
func (c *ObjectCounter[T]) AddMany(values []T) {
	for _, v := range values {
		c.c.Add(c.toId(v))
	}
}
