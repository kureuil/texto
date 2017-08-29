package main

import (
	"github.com/kureuil/texto"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	s := texto.NewServer(":8080", log)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
