package simpleflow

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"strconv"
	"strings"
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

func (s *DeDuplicateSuite) TestDeduplicateObject() {

	// A few values that we can take pointers of
	one := 1
	one_2 := 1
	two := 2
	two_2 := 2

	type Object struct {
		slice   []int
		pointer *int
		value   string
	}

	a := Object{
		slice:   []int{1, 2, 3},
		pointer: &one,
		value:   "one",
	}
	a2 := Object{
		slice:   []int{1, 2, 3},
		pointer: &one_2,
		value:   "one",
	}
	b := Object{
		slice:   []int{2, 3, 4},
		pointer: &two,
		value:   "two",
	}
	b2 := Object{
		slice:   []int{2, 3, 4},
		pointer: &two,
		value:   "two",
	}
	c := Object{
		slice:   []int{2, 3, 4},
		pointer: &two_2,
		value:   "three",
	}
	c2 := Object{
		slice:   []int{4, 5, 6},
		pointer: &two_2,
		value:   "three",
	}

	values := []Object{a, a2, b, b2, c, c2}

	validateSeen := func(dd *ObjectDeduplicator[Object], expectedSeen []Object) {
		for _, expected := range expectedSeen {
			s.True(dd.Seen(expected))
		}
	}

	s.Run("dedupe by value", func() {
		// Create an object deduplicator that looks for unique objects by their "value" field
		dd := NewObjectDeduplicator[Object](func(v Object) string {
			return v.value
		})

		expected := []Object{a, b, c}
		expectedIdx := []int{1, 3, 5}
		s.ElementsMatch(expected, dd.Deduplicate(values))
		dd.Reset() // Need to reset in order to call dd again with the same values
		s.ElementsMatch(expectedIdx, dd.DeduplicateIndices(values))
		validateSeen(dd, expected)
	})

	s.Run("dedupe by pointer value", func() {
		// Create an object deduplicator that looks for unique objects by the value pointed to by the "pointer" field
		dd := NewObjectDeduplicator[Object](func(v Object) string {
			return strconv.Itoa(*v.pointer)
		})

		expected := []Object{a, b}
		expectedIdx := []int{1, 3, 4, 5}
		s.ElementsMatch(expected, dd.Deduplicate(values))
		dd.Reset() // Need to reset in order to call dd again with the same values
		s.ElementsMatch(expectedIdx, dd.DeduplicateIndices(values))
		validateSeen(dd, expected)
	})

	s.Run("dedupe by pointer", func() {
		// Create an object deduplicator that looks for unique objects by the "pointer" field (pointer address)
		dd := NewObjectDeduplicator[Object](func(v Object) string {
			return fmt.Sprintf("%p", v.pointer)
		})

		expected := []Object{a, a2, b, c}
		expectedIdx := []int{3, 5}
		s.ElementsMatch(expected, dd.Deduplicate(values))
		dd.Reset() // Need to reset in order to call dd again with the same values
		s.ElementsMatch(expectedIdx, dd.DeduplicateIndices(values))
		validateSeen(dd, expected)
	})

	s.Run("dedupe by slice value", func() {
		// Create an object deduplicator that looks for unique objects by their "slice" field
		dd := NewObjectDeduplicator[Object](func(v Object) string {
			b := strings.Builder{}
			for _, elem := range v.slice {
				b.WriteString(strconv.Itoa(elem))
			}
			return b.String()
		})

		expected := []Object{a, b, c2}
		expectedIdx := []int{1, 3, 4}
		s.ElementsMatch(expected, dd.Deduplicate(values))
		dd.Reset() // Need to reset in order to call dd again with the same values
		s.ElementsMatch(expectedIdx, dd.DeduplicateIndices(values))
		validateSeen(dd, expected)
	})

	s.Run("dedupe by value and pointer value", func() {
		// Create an object deduplicator that looks for unique objects by multiple fields ("value" and "pointer"'s value)
		dd := NewObjectDeduplicator[Object](func(v Object) string {
			return fmt.Sprintf("%s_%d", v.value, *v.pointer)
		})

		expected := []Object{a, b, c}
		expectedIdx := []int{1, 3, 5}
		s.ElementsMatch(expected, dd.Deduplicate(values))
		dd.Reset() // Need to reset in order to call dd again with the same values
		s.ElementsMatch(expectedIdx, dd.DeduplicateIndices(values))
		validateSeen(dd, expected)
	})

}
