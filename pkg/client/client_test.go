package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/dialer"
	dialerMock "wowpow/internal/pkg/dialer/mock"
	"wowpow/internal/pkg/hash"
	"wowpow/internal/pkg/pow"
	"wowpow/pkg/api/message"
	"wowpow/pkg/client/mock"
)

const testTimeout = time.Minute

func TestGetMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	challengeHashcash, _ := pow.InitHashcash(2, ":0", "secret", hash.NewSHA256())
	challengeBreakErr := fmt.Errorf("challenge break")

	for _, tCase := range []struct {
		name        string
		ctx         func() (context.Context, context.CancelFunc)
		mock        mock.ClientMockParams
		expected    string
		expectedErr error
	}{
		{
			name: "positive close",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mock: mock.ClientMockParams{
				DialerDialTimes: 1,
				DialerDialDoAndReturnFunc: func() func() (dialer.Conn, error) {
					var connFinished bool
					connector := (dialerMock.DialerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnCloseAnyTimes: true,

						ConnReadTimes: 1,
						ConnReadReq:   gomock.Any(),
						ConnReadDoAndReturnFunc: func() func([]byte) (int, error) {
							return func(b []byte) (int, error) {
								defer func() {
									connFinished = true
								}()

								var buf bytes.Buffer
								buf.Write(getMessageBytes(t, message.Message_close, nil, ""))
								buf.Write([]byte{nl})
								return buf.Read(b)
							}
						}(),

						ConnWriteAnyTimes: true,
						ConnWriteReq:      gomock.Any(),
						ConnWriteDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() {
									times++
									if times == 3 {
										// finalyze response
										connFinished = true
									}
								}()

								return len(b), nil
							}
						}(),
					}).NewConn(ctrl)
					var times int
					return func() (dialer.Conn, error) {
						defer func() { times++ }()

						// imitate first connection
						if times == 0 {
							return connector, nil
						}

						// on second iteration wait until connection will be finished
						// return close connection error to finish main routine
						for {
							if connFinished {
								return nil, fmt.Errorf("use of closed network connection")
							}
						}
					}
				}(),
			},
			expectedErr: ErrConnectionClose,
		},
		{
			name: "positive challenge",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mock: mock.ClientMockParams{
				ComputerComputeTimes: 1,
				ComputerComputeReq1:  gomock.Any(),
				ComputerComputeReq2:  gomock.Any(),
				ComputerComputeReq3:  gomock.Any(),
				ComputerComputeRes:   challengeHashcash,

				DialerDialTimes: 1,
				DialerDialDoAndReturnFunc: func() func() (dialer.Conn, error) {
					var connFinished bool
					connector := (dialerMock.DialerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnCloseAnyTimes: true,

						ConnReadTimes: 2,
						ConnReadReq:   gomock.Any(),
						ConnReadDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()

								if times == 0 {
									var buf bytes.Buffer
									// encoded protobuf message
									buf.Write(getMessageBytes(t, message.Message_challenge, challengeHashcash, ""))
									buf.Write([]byte{nl})
									return buf.Read(b)
								}

								// break client on second iteration
								return 0, challengeBreakErr
							}
						}(),

						ConnWriteAnyTimes: true,
						ConnWriteReq:      gomock.Any(),
						ConnWriteDoAndReturnFunc: func() func([]byte) (int, error) {
							return func(b []byte) (int, error) {
								return len(b), nil
							}
						}(),
					}).NewConn(ctrl)
					var times int
					return func() (dialer.Conn, error) {
						defer func() { times++ }()

						// imitate first connection
						if times == 0 {
							return connector, nil
						}

						// on second iteration wait until connection will be finished
						// return close connection error to finish main routine
						for {
							if connFinished {
								return nil, fmt.Errorf("use of closed network connection")
							}
						}
					}
				}(),
			},
			expectedErr: challengeBreakErr,
		},
		{
			name: "positive resource",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mock: mock.ClientMockParams{
				DialerDialTimes: 1,
				DialerDialDoAndReturnFunc: func() func() (dialer.Conn, error) {
					var connFinished bool
					connector := (dialerMock.DialerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnCloseAnyTimes: true,

						ConnReadTimes: 2,
						ConnReadReq:   gomock.Any(),
						ConnReadDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()

								if times == 0 {
									var buf bytes.Buffer
									// encoded protobuf message
									buf.Write(getMessageBytes(t, message.Message_resource, nil, "test"))
									buf.Write([]byte{nl})
									return buf.Read(b)
								}

								// break client on second iteration
								return 0, challengeBreakErr
							}
						}(),

						ConnWriteAnyTimes: true,
						ConnWriteReq:      gomock.Any(),
						ConnWriteDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()

								if times == 1 {
									assert.Equal(t, []byte{nl}, b)
									return 1, nil
								}

								return len(b), nil
							}
						}(),
					}).NewConn(ctrl)
					var times int
					return func() (dialer.Conn, error) {
						defer func() { times++ }()

						// imitate first connection
						if times == 0 {
							return connector, nil
						}

						// on second iteration wait until connection will be finished
						// return close connection error to finish main routine
						for {
							if connFinished {
								return nil, fmt.Errorf("use of closed network connection")
							}
						}
					}
				}(),
			},
			expected: "test",
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a           = assert.New(t)
				computer    = tCase.mock.NewComputer(ctrl)
				dMock       = tCase.mock.NewDialer(ctrl)
				ctx, cancel = tCase.ctx()
			)
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()

			c, _ := NewWoWPoW(computer, dMock)
			actual, err := c.GetMessage(ctx)
			a.Equal(tCase.expected, actual)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func getMessageBytes(t *testing.T, h message.Message_Header, hc *pow.Hashcach, payload string) []byte {
	t.Helper()

	p := &message.Message{
		Header: h,
	}

	if hc != nil {
		p.Response = &message.Message_Hashcach{
			Hashcach: hc.ToProto(),
		}
	}

	if payload != "" {
		p.Response = &message.Message_Payload{
			Payload: payload,
		}
	}

	bin, _ := proto.Marshal(p)

	buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(bin)))
	base64.RawStdEncoding.Encode(buf, bin)

	return buf
}
