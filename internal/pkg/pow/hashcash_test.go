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
		rand     = 10
		date     = time.Unix(1648762844, 0)
		expected = fmt.Sprintf("%s%d%s%d", resource, rand, secret, date.Unix())
		a        = assert.New(t)
		ctrl     = gomock.NewController(t)
		hasher   = mock.NewMockHasher(ctrl)
	)

	hasher.EXPECT().Hash(gomock.Eq(expected)).Return(expected, nil)

	actual, err := extSum(resource, secret, rand, date, hasher)
	a.Nil(err)
	a.Equal(expected, actual)
}

func TestExtSumErr(t *testing.T) {
	var (
		resource    = "resource"
		secret      = "secret"
		rand        = 10
		date        = time.Unix(1648762844, 0)
		expected    = fmt.Sprintf("%s%d%s%d", resource, rand, secret, date.Unix())
		expectedErr = fmt.Errorf("expected error")
		a           = assert.New(t)
		ctrl        = gomock.NewController(t)
		hasher      = mock.NewMockHasher(ctrl)
	)

	hasher.EXPECT().Hash(gomock.Eq(expected)).Return("", expectedErr)

	actual, err := extSum(resource, secret, rand, date, hasher)
	a.Empty(actual)
	a.ErrorIs(err, expectedErr)
}

func TestHashCashMappings(t *testing.T) {
	a := assert.New(t)
	hc := &Hashcach{
		version: versionV1,
		rand:    123,
		counter: 10,
	}

	proto := hc.ToProto()
	a.Equal(
		&message.Hashcach{
			Version: versionV1,
			Date:    timestamppb.New(time.Time{}),
			Rand:    "MTIz",
			Counter: "MTA=",
		},
		proto,
	)

	newHC, err := FromProto(proto)
	a.Nil(err)
	a.Equal(hc, newHC)
}

func TestVerifyExt(t *testing.T) {
	hasherErr := fmt.Errorf("expected error")
	for _, tCase := range []struct {
		name              string
		hashcash          *Hashcach
		secret            string
		hasherExpectedReq string
		hasherRes         string
		hasherErr         error
		expectedErr       error
	}{
		{
			name: "success",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			secret:            "secret",
			hasherExpectedReq: "resource10secret1648762844",
			hasherRes:         "resource10secret1648762844",
			hasherErr:         nil,
			expectedErr:       nil,
		},
		{
			name: "wrong ext",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "wrong",
			},
			secret:            "secret",
			hasherExpectedReq: "resource10secret1648762844",
			hasherRes:         "resource10secret1648762844",
			hasherErr:         nil,
			expectedErr:       ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			secret:            "secret",
			hasherExpectedReq: "resource10secret1648762844",
			hasherRes:         "wrong",
			hasherErr:         nil,
			expectedErr:       ErrExtInvalid,
		},
		{
			name: "wrong hasher response",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			secret:            "secret",
			hasherExpectedReq: "resource10secret1648762844",
			hasherRes:         "wrong",
			hasherErr:         nil,
			expectedErr:       ErrExtInvalid,
		},
		{
			name: "hasher error",
			hashcash: &Hashcach{
				resource: "resource",
				rand:     10,
				date:     time.Unix(1648762844, 0),
				ext:      "resource10secret1648762844",
			},
			secret:            "secret",
			hasherExpectedReq: "resource10secret1648762844",
			hasherRes:         "",
			hasherErr:         hasherErr,
			expectedErr:       hasherErr,
		},
	} {
		t.Run(tCase.name, func(t *testing.T) {
			var (
				a      = assert.New(t)
				ctrl   = gomock.NewController(t)
				hasher = mock.NewMockHasher(ctrl)
			)

			hasher.EXPECT().Hash(gomock.Eq(tCase.hasherExpectedReq)).Return(tCase.hasherRes, tCase.hasherErr)

			err := VerifyExt(tCase.secret, hasher)(tCase.hashcash)
			a.ErrorIs(err, tCase.expectedErr)
		})
	}
}

func TestRandomBytes(t *testing.T) {
	assert.Greater(t, randomBytes(), 0)
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
