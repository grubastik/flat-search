package main

import (
	"flag"
	"github.com/grubastik/flat-search/config"
	"github.com/grubastik/flat-search/db"
	"github.com/grubastik/flat-search/email"
	"github.com/grubastik/flat-search/error"
	"github.com/grubastik/flat-search/aggregator"
)

func main() {
	config := config.MustNewConfig(*flag.String("config", "./config.json", "Path to the config file"))
	storage, err := db.NewDb(config)
	error.DebugError(err)
	defer storage.Close()

	email.NewConnection(config)

	err = aggregator.ProcessAdverts(config.GetSreality())
	error.DebugError(err)
}
