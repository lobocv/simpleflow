package simpleflow

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func generateSeries(n int) (series []int) {
	for ii := 0; ii < n; ii++ {
		series = append(series, ii)
	}
	return
}

type FanSuite struct {
	suite.Suite
}

func TestFan(t *testing.T) {
	s := new(FanSuite)
	suite.Run(t, s)
}

func (s *FanSuite) TestFanOutAndIn() {
	N := 5

	// Generate some data on a channel (source for fan out)
	source := make(chan int, N)
	data := generateSeries(N)
	for _, v := range data {
		source <- v
	}
	close(source)

	// Fan out to two channels. Each must get a copy of the data
	fanoutSink1 := make(chan int, N)
	fanoutSink2 := make(chan int, N)
	FanOutAndClose(source, fanoutSink1, fanoutSink2)

	// Fan them back in to a single channel. We should get the original source data with two copies of each item
	fanInSink := make(chan int, 2*N)
	FanInAndClose(fanInSink, fanoutSink1, fanoutSink2)
	faninResults := ChannelToSlice(fanInSink)

	s.ElementsMatch(faninResults, append(generateSeries(N), generateSeries(N)...))
}
