package simpleflow

import (
	"context"
	"fmt"
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
	f := func(_ context.Context, v int) error {
		out.Set(v, v*v)
		return nil
	}
	errors := WorkerPoolFromSlice(ctx, items, nWorkers, f)

	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
	s.Empty(errors)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromSliceWithErrors() {
	ctx := context.Background()
	items := []int{0, 1, 2, 3, 4, 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, v int) error {
		if v < 3 {
			return fmt.Errorf("%d", v)
		}
		out.Set(v, v)
		return nil
	}
	errors := WorkerPoolFromSlice(ctx, items, nWorkers, f)

	expected := map[int]int{3: 3, 4: 4, 5: 5}
	s.Equal(expected, out.m)
	s.ElementsMatch(errors, []error{
		fmt.Errorf("0"),
		fmt.Errorf("1"),
		fmt.Errorf("2"),
	})
}

func (s *WorkerPoolSuite) TestCancelWorkerPoolFromSlice() {
	ctx, cancel := context.WithCancel(context.Background())
	items := []int{0, 1, 2, 3, 4, 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2

	f := func(_ context.Context, v int) error {
		if v > 2 {
			cancel()
			return nil
		}
		out.Set(v, v*v)
		return nil
	}
	WorkerPoolFromSlice(ctx, items, nWorkers, f)

	// Only keys less than 3 should be processed
	expected := map[int]int{0: 0, 1: 1, 2: 4}
	s.Equal(expected, out.m)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromChan() {
	ctx := context.Background()
	N := 5
	itemChan := make(chan int, N)
	LoadChannel(itemChan, generateSeries(N)...)
	close(itemChan)

	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, v int) error {
		out.Set(v, v*v)
		return nil
	}
	WorkerPoolFromChan(ctx, itemChan, nWorkers, f)
	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16}
	s.Equal(expected, out.m)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromChanCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	N := 100
	itemChan := make(chan int, N)
	LoadChannel(itemChan, generateSeries(N)...)
	close(itemChan)

	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, v int) error {
		if v > 2 {
			cancel()
			return nil
		}
		out.Set(v, v*v)
		return nil
	}
	errors := WorkerPoolFromChan(ctx, itemChan, nWorkers, f)
	// It's not easy to test exactly how many items should get processed due to race conditions,
	// so for now just check that not all items were processed.
	s.NotEqual(len(out.m), N)
	s.Empty(errors)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromChanWithErrors() {
	ctx := context.Background()
	N := 5
	itemChan := make(chan int, N)
	LoadChannel(itemChan, generateSeries(N)...)
	close(itemChan)

	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, v int) error {
		out.Set(v, v)
		if v < 3 {
			return fmt.Errorf("%d", v)
		}
		return nil
	}
	errors := WorkerPoolFromChan(ctx, itemChan, nWorkers, f)
	expected := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	s.Equal(expected, out.m)
	s.ElementsMatch(errors, []error{
		fmt.Errorf("0"),
		fmt.Errorf("1"),
		fmt.Errorf("2"),
	})
}

func (s *WorkerPoolSuite) TestWorkerPoolFromMap() {
	ctx := context.Background()
	items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, k, v int) error {
		out.Set(v, v*v)
		return nil
	}
	errors := WorkerPoolFromMap(ctx, items, nWorkers, f)

	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
	s.Empty(errors)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromMapCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, k, v int) error {
		if k > 2 {
			cancel()
			return nil
		}
		out.Set(v, v*v)
		return nil
	}
	errors := WorkerPoolFromMap(ctx, items, nWorkers, f)
	// It's not easy to test exactly how many items should get processed due to race conditions,
	// so for now just check that not all items were processed.
	s.NotEqual(len(out.m), len(items))
	s.Empty(errors)
}

func (s *WorkerPoolSuite) TestWorkerPoolFromMapWithErrors() {
	ctx := context.Background()
	items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	out := NewSyncMap(map[int]int{})
	nWorkers := 2
	f := func(_ context.Context, k, v int) error {
		out.Set(v, v*v)
		if v < 3 {
			return fmt.Errorf("%d", v)
		}
		return nil
	}
	errors := WorkerPoolFromMap(ctx, items, nWorkers, f)

	expected := map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
	s.Equal(expected, out.m)
	s.ElementsMatch(errors, []error{
		fmt.Errorf("0"),
		fmt.Errorf("1"),
		fmt.Errorf("2"),
	})
}
