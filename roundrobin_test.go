package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RoundRobinSuite struct {
	suite.Suite
}

func TestRoundRobin(t *testing.T) {
	s := new(RoundRobinSuite)
	suite.Run(t, s)
}

func (s *RoundRobinSuite) TestRoundRobin() {
	N := 9
	// Generate some data on a channel (source for fan out)
	source := make(chan int, N)
	LoadChannel(source, generateSeries(N)...)
	close(source)

	// Round robin the data into two channels, each should have half the data
	fanoutSink1 := make(chan int, N)
	fanoutSink2 := make(chan int, N)
	RoundRobin(source, fanoutSink1, fanoutSink2)
	CloseManyWriters(fanoutSink1, fanoutSink2)

	fanout1Data := ChannelToSlice(fanoutSink1)
	fanout2Data := ChannelToSlice(fanoutSink2)

	s.ElementsMatch(fanout1Data, []int{0, 2, 4, 6, 8})
	s.ElementsMatch(fanout2Data, []int{1, 3, 5, 7})
}
