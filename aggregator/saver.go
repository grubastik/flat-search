package aggregator

import (
	"github.com/grubastik/flat-search/db"
	"github.com/grubastik/flat-search/email"
	"github.com/grubastik/flat-search/models"
)

func saveToDb(l []*models.Advert) error {
	dba := new(db.Advert)
	dbl := new(db.Location)
	dbr := new(db.Realtor)
	dbi := new(db.Image)
	dbp := new(db.Property)
	for _, a := range l {
		existInDb := new(models.Advert)
		err := dba.Load("hash_id", a.HashID, existInDb)
		if err != nil && err != db.ErrAdvertDoesNotExist {
			return err
		}
		if existInDb.ID == 0 {
			//record to db
			err = dba.Insert(a)
			if err != nil {
				return err
			}
			a.Location.AdvertID = a.ID
			err = dbl.Insert(a.Location)
			if err != nil {
				return err
			}
			a.Realtor.AdvertID = a.ID
			err = dbr.Insert(a.Realtor)
			if err != nil {
				return err
			}
			if a.Images != nil {
				for _, i := range a.Images {
					i.AdvertID = a.ID
					err = dbi.Insert(i)
					if err != nil {
						return err
					}
				}
			}
			if a.Properties != nil {
				for _, i := range a.Properties {
					i.AdvertID = a.ID
					err = dbp.Insert(i)
					if err != nil {
						return err
					}
				}
			}
			a.IsNew = true
		}
	}
	return nil
}

func sendEmails(l []*models.Advert) error {
	var err error
	var def *email.Definition
	for _, a := range l {
		if a.IsNew {
			//sendemail
			def, err = email.NewEmail(a)
			if err != nil {
				return err
			}
			err = email.C.Send(def)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
