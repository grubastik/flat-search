package sreality

import (
	"github.com/grubastik/flat-search/models"
	"strconv"
)

type urlPartsDefinition struct {
	CategoryMainCb map[string]string
	CategorySubCb  map[string]string
	CategoryTypeCb map[string]string
}

type location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type seo struct {
	CategoryMainCb int64  `json:"category_main_cb"`
	CategorySubCb  int64  `json:"category_sub_cb"`
	CategoryTypeCb int64  `json:"category_type_cb"`
	Locality       string `json:"locality"`
}

// Advert contains info about advert
type Advert struct {
	Locality string    `json:"locality"`
	HashID   int64     `json:"hash_id"`
	Price    float64   `json:"price"`
	Name     string    `json:"name"`
	Seo      *seo      `json:"seo"`
	Gps      *location `json:"gps"`
}

// Adverts contains list of adverts
type Adverts struct {
	Estates *[]Advert `json:"estates"`
}

// Body defines top-level content of the response
type Body struct {
	Embedded *Adverts `json:"_embedded"`
}

// GetAdverts returns list of the advert in reply
func (b *Body) GetAdverts() *[]Advert {
	return b.Embedded.Estates
}

// Count returns qty of the adverts in list
func (b *Body) Count() int {
	return len(*b.GetAdverts())
}

// GetAdvert returns advert according to the requested index from the list
func (b *Body) GetAdvert(index int) *Advert {
	if b.Count() <= index || index < 0 {
		return nil
	}

	return &(*b.GetAdverts())[index]
}

// ConvertToModel builds model based on the data from response
func (a *Advert) ConvertToModel() *models.Advert {
	var (
		aModel = models.NewAdvert()
		lModel = new(models.Location)
	)
	aModel.Locality = a.Locality
	aModel.HashID = a.HashID
	aModel.Price = a.Price
	aModel.Name = a.Name
	lModel.Lat = a.Gps.Lat
	lModel.Lon = a.Gps.Lon
	aModel.Link = a.getLink()
	aModel.Location = lModel
	return aModel
}

func (a *Advert) getLink() string {
	urlParts := &urlPartsDefinition{
		map[string]string{"1": "byt", "2": "dom", "666": "projekt", "3": "pozemka", "4": "komercni", "5": "ostatni"},
		map[string]string{"2": "1%2Bkk", "3": "1%2B1", "4": "2%2Bkk", "5": "2%2B1", "6": "3%2Bkk", "7": "3%2B1", "8": "4%2Bkk", "9": "4%2B1", "10": "5%2Bkk", "11": "5%2B1", "12": "6-a-vice", "16": "atypicky"},
		map[string]string{"1": "prodej", "2": "pronajem", "3": "drazby"}}
	url := conf.GetSreality().URLDetail
	url += urlParts.CategoryTypeCb[strconv.FormatInt(a.Seo.CategoryTypeCb, 10)] + "/"
	url += urlParts.CategoryMainCb[strconv.FormatInt(a.Seo.CategoryMainCb, 10)] + "/"
	url += urlParts.CategorySubCb[strconv.FormatInt(a.Seo.CategorySubCb, 10)] + "/"
	url += a.Seo.Locality + "/" + strconv.FormatInt(a.HashID, 10)
	return url
}
