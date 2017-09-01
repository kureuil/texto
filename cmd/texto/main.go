package main

import (
	"context"
	"os"

	"github.com/kureuil/texto"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	redisAddr := os.Getenv("REDIS_URL")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}
	broker, err := texto.NewRedisBroker(log, redisAddr)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	ctx := context.Background()
	s, err := texto.NewServer(ctx, log, ":" + port, broker)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
