package pow

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"wowpow/internal/pkg/hash"
	"wowpow/internal/pkg/hash/mock"
	"wowpow/pkg/api/message"
)

func TestExtSum(t *testing.T) {
	var (
		resource = "resource"
		secret   = "secret"
		rand     = []byte{10}
		date     = time.Unix(1648762844, 0)
		expected = fmt.Sprintf("%s%s%s%d", resource, rand, secret, date.Unix())
		a        = assert.New(t)
		hasher   = (mock.HasherMockParams{
			HashTimes:  1,
			HashReq:    gomock.Eq(expected),
			HashRes:    expected,
			HashResErr: nil,
		}).NewHasher(gomock.NewController(t))
	)

	actual, err := extSum(resource, secret, rand, date, hasher)
	a.Nil(err)
	a.Equal(expected, actual)
}

func TestExtSumErr(t *testing.T) {
	var (
		resource    = "resource"
		secret      = "secret"
		rand        = []byte{10}
		date        = time.Unix(1648762844, 0)
		expected    = fmt.Sprintf("%s%s%s%d", resource, rand, secret, date.Unix())
		expectedErr = fmt.Errorf("expected error")
		a           = assert.New(t)
		hasher      = (mock.HasherMockParams{
			HashTimes:  1,
			HashReq:    gomock.Eq(expected),
			HashRes:    "",
			HashResErr: expectedErr,
		}).NewHasher(gomock.NewController(t))
	)

	actual, err := extSum(resource, secret, rand, date, hasher)
	a.Empty(actual)
	a.ErrorIs(err, expectedErr)
}

func TestHashCashMappings(t *testing.T) {
	a := assert.New(t)
	hc := &Hashcach{
		version: versionV1,
		rand:    []byte{123},
		counter: 234,
	}

	proto := hc.ToProto()
	a.Equal(
		&message.Hashcach{
			Version: versionV1,
			Date:    timestamppb.New(time.Time{}),
			Rand:    "ew==",
			Counter: "ZWE=",
		},
		proto,
	)

	newHC, err := FromProto(proto)
	a.Nil(err)
	a.Equal(hc, newHC)
}

func TestVerifyExt(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	ctrl := gomock.NewController(t)

	for _, tCase := range []struct {
		name        string
		hashcash    *Hashcach
		secret      string
		hasherMock  mock.HasherMockParams
		expectedErr error
	}{
		{
			name: "positive",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     []byte{10},
				date:     time.Unix(1648762844, 0),
				ext:      "resource\nsecret1648762844",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret1648762844"),
				HashRes:    "resource\nsecret1648762844",
				HashResErr: nil,
			},
			expectedErr: nil,
		},
		{
			name: "wrong ext",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     []byte{10},
				date:     time.Unix(1648762844, 0),
				ext:      "wrong",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret1648762844"),
				HashRes:    "resource\nsecret1648762844",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     []byte{10},
				date:     time.Unix(1648762844, 0),
				ext:      "resource\nsecret1648762844",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret1648762844"),
				HashRes:    "wrong",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     []byte{10},
				date:     time.Unix(1648762844, 0),
				ext:      "resource\nsecret1648762844",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret1648762844"),
				HashRes:    "wrong",
				HashResErr: nil,
			},
			expectedErr: ErrExtInvalid,
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     []byte{10},
				date:     time.Unix(1648762844, 0),
				ext:      "resource\nsecret1648762844",
			},
			secret: "secret",
			hasherMock: mock.HasherMockParams{
				HashTimes:  1,
				HashReq:    gomock.Eq("resource\nsecret1648762844"),
				HashRes:    "",
				HashResErr: hasherErr,
			},
			expectedErr: hasherErr,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				hasher = tCase.hasherMock.NewHasher(ctrl)
			)

			err := VerifyExt(tCase.secret, hasher)(tCase.hashcash)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func TestRandomBytes(t *testing.T) {
	assert.Greater(t, len(randomBytes()), 0)
}

func TestInitHashcash(t *testing.T) {
	var bits int32 = 3
	var resource = "127.0.0.1"
	var secret = "secret"
	a := assert.New(t)
	summer := hash.NewSHA256()

	actual, err := InitHashcash(bits, resource, secret, summer)
	a.Nil(err)
	a.Equal(bits, actual.bits)
	a.Equal(resource, actual.resource)
}
