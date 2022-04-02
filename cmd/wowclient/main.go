package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"wowpow/internal/pkg/config"
	"wowpow/internal/pkg/dialer"
	"wowpow/internal/pkg/hash"
	"wowpow/internal/pkg/pow"
	"wowpow/pkg/client"
)

const maxRetries = 5

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicf("client stopped with error %s", err)
	}

	hasher := hash.NewSHA256()
	powInstance := pow.New(
		hasher,
		pow.WithValidateExtFunc(pow.VerifyExt(cfg.ServerSecret, hasher)),
		pow.WithChallengeExpDuration(cfg.HashcashChallengeExpDuration),
	)

	wowpow, err := client.NewWoWPoW(
		powInstance,
		dialer.New(cfg.Addr),
		client.WithTimeout(cfg.Timeout),
		client.WithMaxIterations(cfg.ClientMaxIterations),
	)
	if err != nil {
		log.Panicf("client init error %s", err)
	}

	var retry int
	for {
		ctx := context.Background()
		msg, err := wowpow.GetMessage(ctx)
		if err != nil {
			retry++
			if retry > maxRetries {
				log.Panicf("client init error %s", err)
			}
			continue
		}

		retry = 0

		fmt.Println(msg)
		<-time.NewTimer(time.Second).C
	}
}
