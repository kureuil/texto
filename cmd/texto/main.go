package main

import (
	"context"

	"github.com/kureuil/texto"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	broker, err := texto.NewRedisBroker(log, "redis:6379")
	if err != nil {
		log.Fatal(err)
	}
	s, err := texto.NewServer(log, ":8080", broker)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
