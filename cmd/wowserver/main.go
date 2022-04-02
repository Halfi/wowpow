package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"wowpow/assets/quotes"
	"wowpow/internal/pkg/app"
	"wowpow/internal/pkg/config"
	"wowpow/internal/pkg/hash"
	"wowpow/internal/pkg/messenger"
	"wowpow/internal/pkg/pow"
	"wowpow/internal/pkg/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicf("application stopped with error %s", err)
	}

	ctx := context.Background()
	hasher := hash.NewSHA256()
	powInstance := pow.New(
		hasher,
		pow.WithValidateExtFunc(pow.VerifyExt(cfg.ServerSecret, hasher)),
		pow.WithChallengeExpDuration(cfg.HashcashChallengeExpDuration),
	)
	a := app.New(time.Minute)

	msger, err := messenger.New(quotes.Quotes)
	if err != nil {
		log.Panicf("messenger init error %s", err)
	}

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Panicf("application stopped with error %s", err)
	}

	serv := server.New(
		lis,
		hasher,
		powInstance,
		msger,
		server.WithListenersLimit(cfg.ServerListenersLimit),
		server.WithBits(cfg.HashcashBits),
		server.WithSecret(cfg.ServerSecret),
		server.WithTimeout(cfg.Timeout),
	)

	a.Register(serv)

	if err := a.Run(ctx); err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		log.Panicf("application stopped with error %s", err)
	}

	log.Print("application successfully stopped")
}
