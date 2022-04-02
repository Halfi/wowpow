package mock

import (
	"net"

	"github.com/golang/mock/gomock"
)

type ServerMockParams struct {
	VerifierVerifyTimes    int
	VerifierVerifyAnyTimes bool
	VerifierVerifyReq1     gomock.Matcher
	VerifierVerifyReq2     gomock.Matcher
	VerifierVerifyReq3     gomock.Matcher
	VerifierVerifyResErr   error

	MessengerMessageTimes    int
	MessengerMessageAnyTimes bool
	MessengerMessageRes      string

	ListenerAcceptTimes           int
	ListenerAcceptAnyTimes        bool
	ListenerAcceptDoAndReturnFunc interface{}

	ListenerCloseTimes    int
	ListenerCloseAnyTimes bool
	ListenerCloseResErr   error

	ListenerAddrTimes    int
	ListenerAddrAnyTimes bool
	ListenerAddrRes      net.Addr

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

func (p ServerMockParams) NewVerifier(ctrl *gomock.Controller) *MockVerifier {
	mock := NewMockVerifier(ctrl)

	callTimes(
		mock.EXPECT().Verify(p.VerifierVerifyReq1, p.VerifierVerifyReq2, p.VerifierVerifyReq3),
		p.VerifierVerifyTimes,
		p.VerifierVerifyAnyTimes,
	).Return(p.VerifierVerifyResErr)

	return mock
}

func (p ServerMockParams) NewMockMessenger(ctrl *gomock.Controller) *MockMessenger {
	mock := NewMockMessenger(ctrl)

	callTimes(mock.EXPECT().GetMessage(), p.MessengerMessageTimes, p.MessengerMessageAnyTimes).Return(p.MessengerMessageRes)

	return mock
}

func (p ServerMockParams) NewMockListener(ctrl *gomock.Controller) *MockListener {
	mock := NewMockListener(ctrl)

	callTimes(mock.EXPECT().Accept(), p.ListenerAcceptTimes, p.ListenerAcceptAnyTimes).DoAndReturn(p.ListenerAcceptDoAndReturnFunc)

	callTimes(mock.EXPECT().Close(), p.ListenerCloseTimes, p.ListenerCloseAnyTimes).Return(p.ListenerCloseResErr)

	callTimes(mock.EXPECT().Addr(), p.ListenerAddrTimes, p.ListenerAddrAnyTimes).Return(p.ListenerAddrRes)

	return mock
}

func (p ServerMockParams) NewConn(ctrl *gomock.Controller) *MockConn {
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
