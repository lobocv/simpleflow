package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DeDuplicateSuite struct {
	suite.Suite
}

func TestDeduplicate(t *testing.T) {
	s := new(DeDuplicateSuite)
	suite.Run(t, s)
}

func (s *DeDuplicateSuite) TestDeduplicateAdd() {

	dd := NewDeduplicator[int]()

	values := []int{1, 2, 3, 3, 4, 5, 6, 6, 6}

	for ii, v := range values {
		isDupe := ii == 3 || ii == 7 || ii == 8
		seen := dd.Seen(v)
		s.Equal(isDupe, seen, "%v  (index %d) is expected to be seen", v, ii)

		added := dd.Add(v)
		s.Equal(isDupe, !added, "%v  (index %d) is expected to be duplicate", v, ii)
	}
}

func (s *DeDuplicateSuite) TestDeduplicateSlice() {
	dd := NewDeduplicator[int]()
	values := []int{1, 2, 3, 3, 4, 5, 6, 6, 6}
	deduped := dd.Deduplicate(values)
	s.ElementsMatch(deduped, []int{1, 2, 3, 4, 5, 6})

	deduped = dd.Deduplicate(values)
	s.Nil(deduped)

	dd.Reset()
	deduped = dd.Deduplicate(values)
	s.ElementsMatch(deduped, []int{1, 2, 3, 4, 5, 6})

	deduped = Deduplicate(values)
	s.ElementsMatch(deduped, []int{1, 2, 3, 4, 5, 6})

	dd.Reset()
	dedupedIndices := dd.DeduplicateIndices(values)
	s.ElementsMatch(dedupedIndices, []int{3, 7, 8})

}
