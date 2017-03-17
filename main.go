
package main

import (
    "flag"
    "github.com/grubastik/flat-search/sreality"
    _ "github.com/grubastik/flat-search/models"
    "github.com/grubastik/flat-search/db"
    "github.com/grubastik/flat-search/email"
    "github.com/grubastik/flat-search/config"
    "github.com/grubastik/flat-search/error"
)



func main() {
    path := flag.String("config", "./config.json", "Path to the config file")
    
    config := config.MustNewConfig(path)
    storage, err := db.NewDb(config)
    error.DebugError(err)
    db.Storage = storage
    defer db.Storage.Close()

    email.Conn = email.NewEmailConnection(config)

    var urlParameters *sreality.UrlParams = sreality.NewSreality(config)
    err = urlParameters.ProcessAdverts()
    error.DebugError(err)
}
