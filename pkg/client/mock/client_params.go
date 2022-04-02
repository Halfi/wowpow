package mock

import (
	"github.com/golang/mock/gomock"

	"wowpow/internal/pkg/pow"
)

type ClientMockParams struct {
	DialerDialTimes           int
	DialerDialAnyTimes        bool
	DialerDialDoAndReturnFunc interface{}

	ComputerComputeTimes    int
	ComputerComputeAnyTimes bool
	ComputerComputeReq1     gomock.Matcher
	ComputerComputeReq2     gomock.Matcher
	ComputerComputeReq3     gomock.Matcher
	ComputerComputeRes      *pow.Hashcach
	ComputerComputeResErr   error
}

func (p ClientMockParams) NewDialer(ctrl *gomock.Controller) *MockDialer {
	mock := NewMockDialer(ctrl)

	callTimes(mock.EXPECT().Dial(), p.DialerDialTimes, p.DialerDialAnyTimes).DoAndReturn(p.DialerDialDoAndReturnFunc)

	return mock
}

func (p ClientMockParams) NewComputer(ctrl *gomock.Controller) *MockComputer {
	mock := NewMockComputer(ctrl)

	callTimes(
		mock.EXPECT().Compute(p.ComputerComputeReq1, p.ComputerComputeReq2, p.ComputerComputeReq3),
		p.ComputerComputeTimes,
		p.ComputerComputeAnyTimes,
	).Return(p.ComputerComputeRes, p.ComputerComputeResErr)

	return mock
}

func callTimes(c *gomock.Call, times int, anyTimes bool) *gomock.Call {
	if anyTimes {
		return c.AnyTimes()
	}

	return c.Times(times)
}
