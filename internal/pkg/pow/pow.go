package pow

import (
	"context"
	"fmt"
	"time"

	"wowpow/internal/pkg/hash"
)

const (
	zero rune = 48 // ASCII code for number zero

	defaultChallengeDuration = 120 * time.Second
)

var (
	ErrMaxIterationsExceeded = fmt.Errorf("max iterations exceeded")
	ErrWrongResource         = fmt.Errorf("wrong resource")
	ErrChallengeExpired      = fmt.Errorf("challenge expired")
	ErrWrongChallenge        = fmt.Errorf("wrong challenge")
)

// POW proof of work class
type POW struct {
	s hash.Hasher

	validateExtFunc      ValidateExtFunc
	challengeExpDuration time.Duration
}

// New constructor
func New(s hash.Hasher, opts ...Options) *POW {
	p := &POW{s: s}

	for i := range opts {
		opts[i](p)
	}

	if p.challengeExpDuration == 0 {
		p.challengeExpDuration = defaultChallengeDuration
	}

	return p
}

// Compute time waster. Do all useless load.
func (p *POW) Compute(ctx context.Context, h *Hashcach, max int) (*Hashcach, error) {
	if max > 0 {
		for h.counter <= max {
			if err := ctx.Err(); err != nil {
				break
			}

			hash, err := p.s.Hash(h.String())
			if err != nil {
				return nil, fmt.Errorf("calculate pow hash sum error: %w", err)
			}

			if isHashCorrect(hash, int(h.bits)) {
				return h, nil
			}

			h.counter++
		}
	}

	return nil, ErrMaxIterationsExceeded
}

// Verify that hashcash correct and provided by server
func (p *POW) Verify(_ context.Context, h *Hashcach, resource string) error {
	if h.resource != resource {
		return ErrWrongResource
	}

	if h.date.Add(p.challengeExpDuration).Before(time.Now()) {
		return ErrChallengeExpired
	}

	hash, err := p.s.Hash(h.String())
	if err != nil {
		return fmt.Errorf("calculate pow hash sum error: %w", err)
	}

	if !isHashCorrect(hash, int(h.bits)) {
		return ErrWrongChallenge
	}

	if p.validateExtFunc != nil {
		err = p.validateExtFunc(h)
		if err != nil {
			return fmt.Errorf("validation extension error: %w", err)
		}
	}

	return nil
}

func isHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}

	for _, ch := range hash[:zerosCount] {
		if ch != zero {
			return false
		}
	}
	return true
}
