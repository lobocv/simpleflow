package concurgo

import (
	"context"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type SyncMap[K comparable, V any] struct {
	sync.Mutex
	m map[K]V
}

func NewSyncMap[K comparable, V any](m map[K]V) *SyncMap[K, V] {
	return &SyncMap[K, V]{m: m}
}

func (m *SyncMap[K, V]) Set(k K, v V) {
	m.Lock()
	m.m[k] = v
	m.Unlock()
}

type WorkerPoolSuite struct {
	suite.Suite
}

func TestWorkerPool(t *testing.T) {
	s := new(WorkerPoolSuite)
	suite.Run(t, s)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromSlice() {
	ctx := context.Background()
	items := []int{0, 1, 2, 3, 4, 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(v int) {
		out.Set(v, v*v)
	}
	WorkerPoolFromSlice(ctx, items, nWorkers, f)

	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromChan() {
	ctx := context.Background()
	items := []int{0, 1, 2, 3, 4, 5}
	itemChan := make(chan int, len(items))
	for _, v := range items {
		itemChan <- v
	}
	close(itemChan)

	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(v int) {
		out.Set(v, v*v)
	}
	WorkerPoolFromChan(ctx, itemChan, nWorkers, f)
	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromMap() {
	ctx := context.Background()
	items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(k, v int) {
		out.Set(v, v*v)
	}
	WorkerPoolFromMap(ctx, items, nWorkers, f)

	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
}
