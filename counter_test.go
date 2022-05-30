package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"strconv"
	"strings"
	"testing"
)

type CounterSuite struct {
	suite.Suite
}

func TestCounter(t *testing.T) {
	s := new(CounterSuite)
	suite.Run(t, s)
}

func (s *CounterSuite) TestCounter() {

	c := NewCounter[int]()

	values := []struct {
		value, currentCount int
	}{
		{1, 1},
		{1, 2},
		{1, 3},
		{2, 1},
		{2, 2},
		{3, 1},
	}

	for _, v := range values {
		count := c.Add(v.value)
		s.Equal(v.currentCount, count)
		s.Equal(v.currentCount, c.Count(v.value))
	}

	c.Reset()
	for _, v := range values {
		s.Equal(0, c.Count(v.value))
	}
	s.Len(c.counts, 0)

	c.AddMany([]int{1, 1, 1, 2, 2, 3})
	s.Equal(3, c.Count(1))
	s.Equal(2, c.Count(2))
	s.Equal(1, c.Count(3))

}

func (s *CounterSuite) TestObjectCounter() {
	// A few values that we can take pointers of
	one := 1
	one_2 := 1
	two := 2
	two_2 := 2
	three := 3

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
		pointer: &three,
		value:   "four",
	}

	values := []Object{a, a2, b, b2, c, c2}

	s.Run("count objects by value field", func() {
		// Create an object deduplicator that looks for unique objects by their "value" field
		dd := NewObjectCounter[Object](func(v Object) string {
			return v.value
		})

		expected := map[string]int{"one": 2, "two": 2, "three": 1, "four": 1}
		dd.AddMany(values)
		for _, v := range values {
			s.Equal(expected[v.value], dd.Count(v))
		}
	})

	s.Run("count objects by pointer field value", func() {
		// Create an object deduplicator that looks for unique objects by their "value" field
		toID := func(v Object) string {
			return strconv.Itoa(*v.pointer)
		}
		dd := NewObjectCounter[Object](toID)

		expected := map[string]int{"1": 2, "2": 3, "3": 1}
		dd.AddMany(values)
		for _, v := range values {
			s.Equal(expected[toID(v)], dd.Count(v))
		}
		dd.Reset()
		for _, v := range values {
			s.Equal(0, dd.Count(v))
		}

	})
	s.Run("count by slice value", func() {
		// Create an object Counter that looks for unique objects by their "slice" field
		toID := func(v Object) string {
			b := strings.Builder{}
			for _, elem := range v.slice {
				b.WriteString(strconv.Itoa(elem))
			}
			return b.String()
		}
		dd := NewObjectCounter[Object](toID)
		expected := map[string]int{"123": 2, "234": 3, "456": 1}
		// Add values individually for test coverage
		for _, v := range values {
			dd.Add(v)
		}
		for _, v := range values {
			s.Equal(expected[toID(v)], dd.Count(v))
		}
	})
}
