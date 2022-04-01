package pow

import (
	"time"
)

type Options func(*POW)

func WithValidateExtFunc(callback ValidateExtFunc) Options {
	return func(pow *POW) {
		pow.validateExtFunc = callback
	}
}
func WithChallengeExpDuration(callback time.Duration) Options {
	return func(pow *POW) {
		pow.challengeExpDuration = callback
	}
}
