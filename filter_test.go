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

func (s *FilterSuite) TestFilterInplace() {

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

func (s *FilterSuite) TestFilter() {

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
