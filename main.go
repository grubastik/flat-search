
package main

import (
    "./sreality"
    _ "./models"
    "./db"
    "./config"
    "./email"
)



func main() {
    config := config.New()

    var urlParameters *sreality.UrlParams = sreality.New(config)
    adverts := urlParameters.MakeRequest()

    db.Storage = db.NewStorage(config)
    defer db.Storage.Close()

    for _, advert := range *(adverts.GetAdverts()) {
        aModel := advert.ConvertToModel()
        if (!aModel.ExistsInDbByHashId()) {
            //record to db
            aModel.Insert()
            //sendemail
            email := email.New(config)
            email.PrepareData(aModel)
            email.Send()
        }
    }
}
