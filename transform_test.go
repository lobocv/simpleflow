package simpleflow

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransformSuite struct {
	suite.Suite
}

func TestTransform(t *testing.T) {
	s := new(TransformSuite)
	suite.Run(t, s)
}

func (s *TransformSuite) TestTransform() {

	out := Transform([]int{1, 2, 3}, func(t int) string {
		return strconv.Itoa(t)
	})

	expected := []string{"1", "2", "3"}
	require.Equal(s.T(), expected, out)
}

func (s *TransformSuite) TestTransformAndFilter() {

	out := TransformAndFilter([]int{1, 2, 3, 4, 5}, func(t int) (int, bool) {
		return 2 * t, t%2 == 0
	})

	expected := []int{4, 8}
	require.Equal(s.T(), expected, out)
}
