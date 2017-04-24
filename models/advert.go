package models

import (
	"time"
)

const statusNew = "new"

// Advert defines which advert info will be stored/processed
type Advert struct {
	ID          int64
	Locality    string
	Link        string
	HashID      int64
	Price       float64
	Name        string
	Status      string
	Created     int64
	Location    *Location
	IsNew       bool
	Description string
	Properties  []Property
	Images      []Image
	Realtor     Realtor
}

// Location defines info about location of the advert
type Location struct {
	ID       int64
	AdvertID int64
	Lat      float64
	Lon      float64
}

// Property defines info about flat property
type Property struct {
	ID       int64
	AdvertID int64
	Name     string
	Value    string
}

// Image defines url of the image
type Image struct {
	ID       int64
	AdvertID int64
	URL      string
}

// Realtor defines basic info about realtor
type Realtor struct {
	ID           int64
	AdvertID     int64
	Name         string
	Phone        string
	Email        string
	Company      string
	CompanyPhone string
	CompanyICO   string
}

// NewAdvert creates new structure for advert
func NewAdvert() *Advert {
	var model = new(Advert)
	model.SetCreatedAt()
	model.SetStatusInitial()
	return model
}

// SetCreatedAt set created to current time
func (a *Advert) SetCreatedAt() {
	a.Created = time.Now().Unix()
}

// SetStatusInitial set status to new
func (a *Advert) SetStatusInitial() {
	a.Status = statusNew
}

// GetCreated converts timestamp to time struct and returns it
func (a *Advert) GetCreated() time.Time {
	return time.Unix(a.Created, 0)
}

