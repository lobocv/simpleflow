package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type SegmentSuite struct {
	suite.Suite
}

func TestSegment(t *testing.T) {
	s := new(SegmentSuite)
	suite.Run(t, s)
}

// Segment the items into two slices, one with even values and one with odd values
func (s *SegmentSuite) TestSegmentSlice() {
	items := []int{0, 1, 2, 3, 4, 5}

	segments := SegmentSlice(items, func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	s.Len(segments, 2)
	s.ElementsMatch(segments["even"], []int{0, 2, 4})
	s.ElementsMatch(segments["odd"], []int{1, 3, 5})
}

// Segment the items into two maps, one with capitalized keys and one with lower case keys
func (s *SegmentSuite) TestSegmentMap() {
	items := map[string]int{"One": 1, "Two": 2, "three": 3, "four": 4}

	segments := SegmentMap(items, func(k string, v int) string {
		if k == strings.ToLower(k) {
			return "lowercase"
		}
		return "capitalized"
	})
	s.Len(segments, 2)
	s.Equal(segments["capitalized"], map[string]int{"One": 1, "Two": 2})
	s.Equal(segments["lowercase"], map[string]int{"three": 3, "four": 4})
}
