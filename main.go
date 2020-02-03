package main

import (
	"flag"
	"log"

	"github.com/grubastik/flat-search/aggregator"
	"github.com/grubastik/flat-search/config"
	"github.com/grubastik/flat-search/db"
	"github.com/grubastik/flat-search/email"
)

func main() {
	log.Println("start")
	cfg := config.MustNewConfig(*flag.String("config", "./config.json", "Path to the config file"))
	storage, err := db.NewDb(cfg)
	if err != nil {
		log.Fatal(err, "can't open DB")
	}
	defer storage.Close()

	err = storage.Db.Ping()
	if err != nil {
		log.Fatal(err, "ping error")
	}

	email.NewConnection(cfg)

	err = aggregator.ProcessAdverts(cfg.GetSreality())
	if err != nil {
		log.Fatal(err, "process adverts error")
	}
	log.Println("end")
}
