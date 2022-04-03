package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/hash"
	hashMock "wowpow/internal/pkg/hash/mock"
	"wowpow/internal/pkg/pow"
	"wowpow/internal/pkg/server/mock"
	"wowpow/pkg/api/message"
)

const testTimeout = time.Minute

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	for _, tCase := range []struct {
		name        string
		ctx         func() (context.Context, context.CancelFunc)
		mocks       mock.ServerMockParams
		hasherMock  hashMock.HasherMockParams
		expectedErr error
	}{
		{
			name: "positive close",
			ctx: func() (context.Context, context.CancelFunc) {
				// add time limitation
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mocks: mock.ServerMockParams{
				ListenerAddrTimes: 1,
				ListenerAddrRes:   new(net.TCPAddr),

				ListenerCloseTimes: 1,

				ListenerAcceptTimes: 2,
				ListenerAcceptDoAndReturnFunc: func() func() (Conn, error) {
					var connFinished bool
					connector := (mock.ServerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnRemoteAddrTimes: 1,
						ConnRemoteAddrRes:   new(net.TCPAddr),

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
					}).NewConn(ctrl)
					var times int
					return func() (Conn, error) {
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
			hasherMock: hashMock.HasherMockParams{
				HashAnyTimes: true,
				HashReq:      gomock.Any(),
			},
			expectedErr: nil,
		},
		{
			name: "positive challenge",
			ctx: func() (context.Context, context.CancelFunc) {
				// add time limitation
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mocks: mock.ServerMockParams{
				ListenerAddrTimes: 1,
				ListenerAddrRes:   new(net.TCPAddr),

				ListenerCloseTimes: 1,

				ListenerAcceptTimes: 2,
				ListenerAcceptDoAndReturnFunc: func() func() (Conn, error) {
					var connFinished bool
					connector := (mock.ServerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnRemoteAddrTimes: 1,
						ConnRemoteAddrRes:   new(net.TCPAddr),

						ConnCloseAnyTimes: true,

						ConnReadTimes: 2,
						ConnReadReq:   gomock.Any(),
						ConnReadDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()
								if times == 0 {
									// first time
									var buf bytes.Buffer
									buf.Write(getMessageBytes(t, message.Message_challenge, nil, ""))
									buf.Write([]byte{nl})
									return buf.Read(b)
								}

								// return end of file error to close connection reading routine
								return 0, io.EOF
							}
						}(),

						ConnWriteTimes: 2,
						ConnWriteReq:   gomock.Any(),
						ConnWriteDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()

								if times == 0 {
									assert.Greater(t, len(b), 0)
									return len(b), nil
								}

								// finalyze response
								connFinished = true
								assert.Equal(t, []byte{nl}, b)
								return 1, nil
							}
						}(),
					}).NewConn(ctrl)
					var times int
					return func() (Conn, error) {
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
			hasherMock: hashMock.HasherMockParams{
				HashAnyTimes: true,
				HashReq:      gomock.Any(),
			},
			expectedErr: nil,
		},
		{
			name: "positive resource",
			ctx: func() (context.Context, context.CancelFunc) {
				// add time limitation
				return context.WithTimeout(context.Background(), testTimeout)
			},
			mocks: mock.ServerMockParams{
				ListenerAddrTimes: 1,
				ListenerAddrRes:   new(net.TCPAddr),

				ListenerCloseTimes: 1,

				MessengerMessageTimes: 1,
				MessengerMessageRes:   "",

				VerifierVerifyAnyTimes: true,
				VerifierVerifyReq1:     gomock.Any(),
				VerifierVerifyReq2:     gomock.Any(),
				VerifierVerifyReq3:     gomock.Any(),

				ListenerAcceptTimes: 2,
				ListenerAcceptDoAndReturnFunc: func() func() (Conn, error) {
					var connFinished bool
					connector := (mock.ServerMockParams{
						ConnSetDeadlineTimes: 1,
						ConnSetDeadlineReq:   gomock.Any(),

						ConnRemoteAddrTimes: 1,
						ConnRemoteAddrRes:   new(net.TCPAddr),

						ConnCloseAnyTimes: true,

						ConnReadTimes: 2,
						ConnReadReq:   gomock.Any(),
						ConnReadDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()
								// first time
								if times == 0 {
									var buf bytes.Buffer

									challengeHashcash, _ := pow.InitHashcash(2, ":0", "secret", hash.NewSHA256())

									// encoded protobuf message
									buf.Write(getMessageBytes(t, message.Message_resource, challengeHashcash, ""))
									buf.Write([]byte{nl})
									return buf.Read(b)
								}

								// return end of file error to close connection reading routine
								return 0, io.EOF
							}
						}(),

						ConnWriteTimes: 2,
						ConnWriteReq:   gomock.Any(),
						ConnWriteDoAndReturnFunc: func() func([]byte) (int, error) {
							var times int
							return func(b []byte) (int, error) {
								defer func() { times++ }()

								if times == 0 {
									assert.Greater(t, len(b), 0)
									return len(b), nil
								}

								// finalyze response
								connFinished = true
								assert.Equal(t, []byte{nl}, b)
								return 1, nil
							}
						}(),
					}).NewConn(ctrl)
					var times int
					return func() (Conn, error) {
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
			hasherMock: hashMock.HasherMockParams{
				HashAnyTimes: true,
				HashReq:      gomock.Any(),
			},
			expectedErr: nil,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a         = assert.New(t)
				listener  = tCase.mocks.NewMockListener(ctrl)
				hasher    = tCase.hasherMock.NewHasher(ctrl)
				verifier  = tCase.mocks.NewVerifier(ctrl)
				messenger = tCase.mocks.NewMockMessenger(ctrl)
			)

			ctx, cancel := tCase.ctx()
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()

			err := New(listener, hasher, verifier, messenger).Run(ctx)

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
