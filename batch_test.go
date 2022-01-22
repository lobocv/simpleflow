package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type BatchSuite struct {
	suite.Suite
}

func TestBatch(t *testing.T) {
	s := new(BatchSuite)
	suite.Run(t, s)
}

func (s *BatchSuite) TestBatchSlice() {
	items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s.Run("integer number of batches", func() {
		batches := BatchSlice(items, 2)
		expected := [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}
		s.Equal(expected, batches)
	})

	s.Run("fractional number of batches", func() {
		batches := BatchSlice(items, 6)
		expected := [][]int{{0, 1, 2, 3, 4, 5}, {6, 7, 8, 9}}
		s.Equal(expected, batches)
	})

}

func (s *BatchSuite) TestBatchChan() {

	// creates a channel populated with values
	initData := func() chan int {
		N := 10
		items := make(chan int, N)
		LoadChannel(items, generateSeries(N)...)
		close(items)
		return items
	}

	s.Run("integer number of batches", func() {
		items := initData()
		batchsize := 2
		expectedBatches := 5
		out := make(chan []int, expectedBatches)
		BatchChan(items, batchsize, out)
		close(out)
		batches := ChannelToSlice(out)
		expected := [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}
		s.Equal(expected, batches)
	})

	s.Run("fractional number of batches", func() {
		items := initData()
		batchsize := 6
		expectedBatches := 2
		out := make(chan []int, expectedBatches)
		BatchChan(items, batchsize, out)
		close(out)
		batches := ChannelToSlice(out)
		expected := [][]int{{0, 1, 2, 3, 4, 5}, {6, 7, 8, 9}}
		s.Equal(expected, batches)
	})

}

func (s *BatchSuite) TestMapSlice() {
	items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9}

	s.Run("integer number of batches", func() {
		batches := BatchMap(items, 2)
		s.Len(batches, 5)
		for ii := 0; ii < 5; ii++ {
			s.Len(batches[ii], 2)
		}
	})

	s.Run("fractional number of batches", func() {
		batches := BatchMap(items, 6)
		s.Len(batches, 2)
		s.Len(batches[0], 6)
		s.Len(batches[1], 4)
	})

}
