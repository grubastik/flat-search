package main

import (
	"flag"
	"github.com/grubastik/flat-search/config"
	"github.com/grubastik/flat-search/db"
	"github.com/grubastik/flat-search/email"
	"github.com/grubastik/flat-search/error"
	_ "github.com/grubastik/flat-search/models"
	"github.com/grubastik/flat-search/sreality"
)

func main() {
	config := config.MustNewConfig(flag.String("config", "./config.json", "Path to the config file"))
	storage, err := db.NewDb(config)
	error.DebugError(err)
	defer storage.Close()

	email.NewConnection(config)

	urlParameters := sreality.NewSreality(config)
	err = urlParameters.ProcessAdverts()
	error.DebugError(err)
}
