package pow

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"wowpow/internal/pkg/hash/mock"
)

func TestIsHashCorrect(t *testing.T) {
	for _, tCase := range []struct {
		name     string
		hash     string
		zeros    int
		expected bool
	}{
		{
			name:     "success",
			hash:     "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			zeros:    5,
			expected: true,
		},
		{
			name:     "wrong 5 zeros",
			hash:     "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			zeros:    6,
			expected: false,
		},
		{
			name:     "wrong 0",
			hash:     "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
			zeros:    5,
			expected: false,
		},
		{
			name:     "too short",
			hash:     "0000",
			zeros:    6,
			expected: false,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			actual := isHashCorrect(tCase.hash, tCase.zeros)
			assert.Equal(t, tCase.expected, actual)
		})
	}
}

func TestPowCompute(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	deadCTX, deadCancel := context.WithCancel(context.Background())
	deadCancel()

	for _, tCase := range []struct {
		name              string
		ctx               context.Context
		hashcash          *Hashcach
		max               int
		hasherExpectedReq gomock.Matcher
		hasherCallTimes   int
		hasherRes         string
		hasherErr         error
		expected          *Hashcach
		expectedErr       error
	}{
		{
			name: "success",
			ctx:  ctx,
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			max:               1,
			hasherExpectedReq: gomock.Eq("0:5:1648762844:resource:resource10secret1648762844:MTA=:MA=="),
			hasherCallTimes:   1,
			hasherRes:         "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			hasherErr:         nil,
			expected: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
				counter:  0,
			},
			expectedErr: nil,
		},
		{
			name: "hasher error",
			ctx:  ctx,
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			max:               1,
			hasherCallTimes:   1,
			hasherExpectedReq: gomock.Eq("0:0:1648762844:resource:resource10secret1648762844:MTA=:MA=="),
			hasherRes:         "",
			hasherErr:         hasherErr,
			expected:          nil,
			expectedErr:       hasherErr,
		},
		{
			name: "deadline exceeded",
			ctx:  ctx,
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			max:               1,
			hasherCallTimes:   2,
			hasherExpectedReq: gomock.Any(),
			hasherRes:         "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
			hasherErr:         nil,
			expected:          nil,
			expectedErr:       ErrMaxIterationsExceeded,
		},
		{
			name: "dead ctx",
			ctx:  deadCTX,
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			max:               1,
			hasherCallTimes:   0,
			hasherExpectedReq: gomock.Any(),
			expected:          nil,
			expectedErr:       ErrMaxIterationsExceeded,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				ctrl   = gomock.NewController(t)
				hasher = mock.NewMockHasher(ctrl)
				pow    = New(hasher)
			)

			hasher.EXPECT().Hash(tCase.hasherExpectedReq).Times(tCase.hasherCallTimes).Return(tCase.hasherRes, tCase.hasherErr)

			actual, err := pow.Compute(tCase.ctx, tCase.hashcash, tCase.max)
			a.Equal(tCase.expected, actual)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func TestPowVerify(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, tCase := range []struct {
		name              string
		powOptions        []Options
		hashcash          *Hashcach
		resource          string
		hasherExpectedReq gomock.Matcher
		hasherCallTimes   int
		hasherRes         string
		hasherErr         error
		expectedErr       error
	}{
		{
			name: "success",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     now,
				ext:      "resource10secret1648762844",
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			hasherErr:         nil,
			expectedErr:       nil,
		},
		{
			name: "success validate ext",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     now,
				ext:      "resource10secret1648762844",
			},
			powOptions: []Options{
				WithValidateExtFunc(func(h *Hashcach) error {
					assert.Equal(
						t,
						&Hashcach{
							bits:     5,
							resource: "resource",
							rand:     10,
							date:     now,
							ext:      "resource10secret1648762844",
						},
						h,
					)
					return nil
				}),
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			hasherErr:         nil,
			expectedErr:       nil,
		},
		{
			name: "success duration",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     now.Add(50 * time.Second),
				ext:      "resource10secret1648762844",
			},
			powOptions: []Options{
				WithChallengeExpDuration(time.Minute),
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			hasherErr:         nil,
			expectedErr:       nil,
		},
		{
			name: "wrong resource",
			hashcash: &Hashcach{
				resource: "resource",
			},
			resource:          "resource2",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   0,
			expectedErr:       ErrWrongResource,
		},
		{
			name: "challenge expired",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   0,
			expectedErr:       ErrChallengeExpired,
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Now(),
				ext:      "resource10secret1648762844",
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "",
			hasherErr:         hasherErr,
			expectedErr:       hasherErr,
		},
		{
			name: "wrong hash",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Now(),
				ext:      "resource10secret1648762844",
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "d59d15c9a1842bc4563897803799e94f1f242d7e7e8c618f047e068211543998",
			hasherErr:         nil,
			expectedErr:       ErrWrongChallenge,
		},
		{
			name: "validate ext error",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     time.Now(),
				ext:      "resource10secret1648762844",
			},
			powOptions: []Options{
				WithValidateExtFunc(func(h *Hashcach) error {
					return hasherErr
				}),
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   1,
			hasherRes:         "00000e89df98a05e524fdcd29d8040d64d0259e2d5109ca1998e567a3c1c1c68",
			hasherErr:         nil,
			expectedErr:       hasherErr,
		},
		{
			name: "error duration",
			hashcash: &Hashcach{
				bits:     5,
				resource: "resource",
				rand:     10,
				date:     now.Add(-2 * time.Minute),
				ext:      "resource10secret1648762844",
			},
			powOptions: []Options{
				WithChallengeExpDuration(time.Minute),
			},
			resource:          "resource",
			hasherExpectedReq: gomock.Any(),
			hasherCallTimes:   0,
			expectedErr:       ErrChallengeExpired,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				ctrl   = gomock.NewController(t)
				hasher = mock.NewMockHasher(ctrl)
				pow    = New(hasher, tCase.powOptions...)
			)

			hasher.EXPECT().Hash(tCase.hasherExpectedReq).Times(tCase.hasherCallTimes).Return(tCase.hasherRes, tCase.hasherErr)

			err := pow.Verify(ctx, tCase.hashcash, tCase.resource)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}
