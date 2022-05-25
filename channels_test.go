package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChannelsSuite struct {
	suite.Suite
}

func TestChannels(t *testing.T) {
	s := new(ChannelsSuite)
	suite.Run(t, s)
}

func (s *ChannelsSuite) TestCloseMany() {
	ch1 := make(chan int, 3)
	ch1 <- 1
	ch2 := make(chan int, 3)
	ch2 <- 2
	CloseMany(ch1, ch2)

	var values []int
	values = ChannelIntoSlice(ch1, values)
	values = ChannelIntoSlice(ch2, values)
	s.ElementsMatch([]int{1, 2}, values)
}
