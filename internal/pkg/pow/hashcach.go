package pow

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"wowpow/internal/pkg/hash"
	"wowpow/pkg/api/message"
)

const (
	versionV1 = 1
)

var (
	ErrExtInvalid    = fmt.Errorf("extension sum invalid")
	ErrHashcashEmpty = fmt.Errorf("hashcash empty")
)

type ValidateExtFunc func(h *Hashcach) error

// Hashcach struct to marshal and unmarshal hashcach to string or proto buf
type Hashcach struct {
	version  int32
	bits     int32
	date     time.Time
	resource string
	ext      string
	rand     int
	counter  int
}

func newHashcach(
	version int32,
	bits int32,
	date time.Time,
	resource string,
	ext string,
	rand int,
	counter int,
) *Hashcach {
	return &Hashcach{
		version:  version,
		bits:     bits,
		date:     date,
		resource: resource,
		ext:      ext,
		rand:     rand,
		counter:  counter,
	}
}

// FromProto returns Hashcach struct from proto message
func FromProto(hashcach *message.Hashcach) (*Hashcach, error) {
	if hashcach == nil {
		return nil, ErrHashcashEmpty
	}

	var t time.Time
	if d := hashcach.GetDate(); d != nil {
		t = d.AsTime()
	}

	counterDecoded, err := base64.StdEncoding.DecodeString(hashcach.GetCounter())
	if err != nil {
		return nil, fmt.Errorf("counter base64 decode error: %w", err)
	}

	counter, err := strconv.Atoi(string(counterDecoded))
	if err != nil {
		return nil, fmt.Errorf("counter parse int error: %w", err)
	}

	randDecoded, err := base64.StdEncoding.DecodeString(hashcach.GetRand())
	if err != nil {
		return nil, fmt.Errorf("rand base64 decode error: %w", err)
	}

	rand, err := strconv.Atoi(string(randDecoded))
	if err != nil {
		return nil, fmt.Errorf("rand parse int error: %w", err)
	}

	return newHashcach(
		versionV1,
		hashcach.GetBits(),
		t,
		hashcach.GetResource(),
		hashcach.GetExt(),
		rand,
		counter,
	), nil
}

// InitHashcash initiate new hashcash
func InitHashcash(bits int32, resource, secret string, hasher hash.Hasher) (*Hashcach, error) {
	t := time.Now()
	randBytes := randomBytes()

	extSum, err := extSum(resource, secret, randBytes, t, hasher)
	if err != nil {
		return nil, fmt.Errorf("calculate hashcash ext hash sum error: %w", err)
	}

	return newHashcach(
		versionV1,
		bits,
		t,
		resource,
		extSum,
		randBytes,
		0,
	), nil
}

// String implements fmt.Stringer interface to get string hashcash
func (h *Hashcach) String() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(int(h.version)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(h.bits)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(h.date.Unix())))
	buf.WriteString(":")
	buf.WriteString(h.resource)
	buf.WriteString(":")
	buf.WriteString(h.ext)
	buf.WriteString(":")
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.rand))))
	buf.WriteString(":")
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))))
	return buf.String()
}

// ToProto map hashcash to proto struct
func (h *Hashcach) ToProto() *message.Hashcach {
	return &message.Hashcach{
		Version:  h.version,
		Bits:     h.bits,
		Date:     timestamppb.New(h.date),
		Resource: h.resource,
		Ext:      h.ext,
		Rand:     base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.rand))),
		Counter:  base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
	}
}

func randomBytes() int {
	b, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		b = big.NewInt(mrand.Int63n(math.MaxInt64))
	}

	return int(b.Int64())
}

// extSum generates hash sum with hasher interface from fields:
//    - resource  - ip address
//    - randBytes - random number
//    - secret    - secret known only on server
//    - time      - timestamp
func extSum(resource, secret string, randBytes int, t time.Time, hasher hash.Hasher) (string, error) {
	var ext bytes.Buffer
	ext.WriteString(resource)
	ext.WriteString(strconv.Itoa(randBytes))
	ext.WriteString(secret)
	ext.WriteString(strconv.Itoa(int(t.Unix())))

	extSum, err := hasher.Hash(ext.String())
	if err != nil {
		return "", fmt.Errorf("calculate hashcash ext hash sum error: %w", err)
	}

	return extSum, nil
}

// VerifyExt verify extension from hashcash to validate hashcash was provided by server.
// See extSum description for hash generating details
func VerifyExt(secret string, hasher hash.Hasher) ValidateExtFunc {
	return func(h *Hashcach) error {
		extSum, err := extSum(h.resource, secret, h.rand, h.date, hasher)
		if err != nil {
			return fmt.Errorf("verify ext sum error: %w", err)
		}

		if h.ext != extSum {
			return ErrExtInvalid
		}

		return nil
	}
}
