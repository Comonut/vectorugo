package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/Comonut/vectorugo/server"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	dimension := flag.Int("dim", -1, "vector dimension value")
	name := flag.String("name", "vectors", "index name, default is vectors")
	persitance := flag.Bool("persistance", false, "weather to use in memory storage or persistant storage")
	flag.Parse()
	if *persitance {

		if *dimension < 0 {
			logrus.Error("Invalid dimension value set a valid with -dim flag ex. : -dim 256")
			return
		}
	}
	logrus.Info("Starting server")
	server.Init(uint32(*dimension), *name, *persitance)
}
