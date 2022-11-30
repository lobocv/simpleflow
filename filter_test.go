package simpleflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FilterSuite struct {
	suite.Suite
}

func TestFilter(t *testing.T) {
	s := new(FilterSuite)
	suite.Run(t, s)
}

func (s *FilterSuite) TestFilterSliceInplace() {

	testCases := map[string]struct {
		in       []int
		expected []int
	}{
		"ordered":   {in: []int{-1, 0, 1, 2, 3}, expected: []int{1, 2, 3}},
		"unordered": {in: []int{5, -2, 3, 1, 0, -3, -5, -6}, expected: []int{5, 3, 1}},
	}

	getPositive := func(t int) bool {
		return t > 0
	}

	for name, tc := range testCases {
		out := FilterSliceInplace(tc.in, getPositive)

		require.Equal(s.T(), tc.expected, out, "failed test case", name)
		// Check that capacity has not changed (not allocating a new slice)
		require.Equal(s.T(), len(tc.in), cap(out))
		// Check that the memory address of the pointer has not changed
		require.Equal(s.T(), fmt.Sprintf("%p", tc.in), fmt.Sprintf("%p", out))
	}

}

func (s *FilterSuite) TestFilterSlice() {

	testCases := map[string]struct {
		in       []int
		expected []int
	}{
		"ordered":   {in: []int{-1, 0, 1, 2, 3}, expected: []int{1, 2, 3}},
		"unordered": {in: []int{5, -2, 3, 1, 0, -3, -5, -6}, expected: []int{5, 3, 1}},
	}

	getPositive := func(t int) bool {
		return t > 0
	}

	for name, tc := range testCases {
		out := FilterSlice(tc.in, getPositive)

		require.Equal(s.T(), tc.expected, out, "failed test case: '%s'", name)
		// Check that the memory address of the pointer has not changed
		require.NotEqual(s.T(), fmt.Sprintf("%p", tc.in), fmt.Sprintf("%p", out))
	}

}

func (s *FilterSuite) TestFilterMap() {

	in := map[string]int{"negative_one": -1, "zero": 0, "one": 1, "two": 2, "three": 3}
	expected := map[string]int{"one": 1, "two": 2, "three": 3}

	getPositive := func(key string, value int) bool {
		return value > 0
	}

	out := FilterMap(in, getPositive)

	require.Equal(s.T(), expected, out)

}

func (s *FilterSuite) TestFilterMapInplace() {

	in := map[string]int{"negative_one": -1, "zero": 0, "one": 1, "two": 2, "three": 3}
	expected := map[string]int{"one": 1, "two": 2, "three": 3}

	getPositive := func(key string, value int) bool {
		return value > 0
	}

	FilterMapInplace(in, getPositive)

	require.Equal(s.T(), expected, in)
}
