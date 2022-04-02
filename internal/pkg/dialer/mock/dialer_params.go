package mock

import (
	"net"

	"github.com/golang/mock/gomock"
)

type DialerMockParams struct {
	ConnReadTimes           int
	ConnReadAnyTimes        bool
	ConnReadReq             gomock.Matcher
	ConnReadDoAndReturnFunc interface{}

	ConnWriteTimes           int
	ConnWriteAnyTimes        bool
	ConnWriteReq             gomock.Matcher
	ConnWriteDoAndReturnFunc interface{}

	ConnCloseTimes    int
	ConnCloseAnyTimes bool
	ConnCloseResErr   error

	ConnRemoteAddrTimes    int
	ConnRemoteAddrAnyTimes bool
	ConnRemoteAddrRes      net.Addr

	ConnSetDeadlineTimes    int
	ConnSetDeadlineAnyTimes bool
	ConnSetDeadlineReq      gomock.Matcher
	ConnSetDeadlineResErr   error
}

func (p DialerMockParams) NewConn(ctrl *gomock.Controller) *MockConn {
	mock := NewMockConn(ctrl)

	callTimes(mock.EXPECT().Read(p.ConnReadReq), p.ConnReadTimes, p.ConnReadAnyTimes).DoAndReturn(p.ConnReadDoAndReturnFunc)

	callTimes(mock.EXPECT().Write(p.ConnWriteReq), p.ConnWriteTimes, p.ConnWriteAnyTimes).DoAndReturn(p.ConnWriteDoAndReturnFunc)

	callTimes(mock.EXPECT().Close(), p.ConnCloseTimes, p.ConnCloseAnyTimes).Return(p.ConnCloseResErr)

	callTimes(mock.EXPECT().RemoteAddr(), p.ConnRemoteAddrTimes, p.ConnRemoteAddrAnyTimes).Return(p.ConnRemoteAddrRes)

	callTimes(mock.EXPECT().SetDeadline(p.ConnSetDeadlineReq), p.ConnSetDeadlineTimes, p.ConnSetDeadlineAnyTimes).Return(p.ConnSetDeadlineResErr)

	return mock
}

func callTimes(c *gomock.Call, times int, anyTimes bool) *gomock.Call {
	if anyTimes {
		return c.AnyTimes()
	}

	return c.Times(times)
}
