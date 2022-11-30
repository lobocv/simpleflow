package simpleflow

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExtractSuite struct {
	suite.Suite
}

func TestExtract(t *testing.T) {
	s := new(ExtractSuite)
	suite.Run(t, s)
}

func (s *ExtractSuite) TestExtractToSlice() {

	type Object struct {
		Name string
	}

	in := []Object{
		{Name: "John"},
		{Name: "Paul"},
		{Name: "George"},
		{Name: "Ringo"},
		{Name: "Bob"},
	}
	var names []string

	fn := func(t Object) (string, bool) {
		if t.Name == "Bob" {
			return "", false
		}
		return t.Name, true
	}

	names = ExtractToSlice(in, fn, names)

	expected := []string{"John", "Paul", "George", "Ringo"}
	require.Equal(s.T(), expected, names)
}

func (s *ExtractSuite) TestExtractToChannel() {

	type Object struct {
		Name string
	}

	in := []Object{
		{Name: "John"},
		{Name: "Paul"},
		{Name: "George"},
		{Name: "Ringo"},
		{Name: "Bob"},
	}
	names := make(chan string, len(in))

	fn := func(t Object) (string, bool) {
		if t.Name == "Bob" {
			return "", false
		}
		return t.Name, true
	}

	ExtractToChannel(in, fn, names)
	close(names)

	require.Len(s.T(), names, 4)
	
	expected := []string{"John", "Paul", "George", "Ringo"}
	require.Equal(s.T(), expected, ChannelToSlice(names))
}
