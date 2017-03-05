
package sreality

import(
    "./../models"
    "strconv"
)

const domain = "https://www.sreality.cz/detail/";

type urlPartsDefinition struct {
    Category_main_cb map[string]string
    Category_sub_cb map[string]string
    Category_type_cb map[string]string
}

type location struct {
    Id int64
    AdvertId int64
    Lat float64
    Lon float64
}

type seo struct {
    Category_main_cb int64
    Category_sub_cb int64
    Category_type_cb int64
    Locality string
}

type advert struct {
    Id int64
    Locality string
    Hash_id int64
    Price float64
    Name string
    Seo *seo
    Gps *location
}

type adverts struct {
    Estates *[]advert
}

type Body struct {
    Embedded *adverts `json:"_embedded"`
}

func (b *Body) GetAdverts() *[]advert {
    return b.Embedded.Estates;
}

func (b *Body) Count() int {
    return len(*b.GetAdverts())
}

func (b *Body) GetAdvert(index int) *advert {
    if b.Count() <= index || index < 0 {
        return nil
    }
    
    return &(*b.GetAdverts())[index]
}

func (a *advert) ConvertToModel() (*models.Advert){
    var aModel = new(models.Advert)
    var lModel = new(models.Location)
    aModel.SetLocality(a.Locality)
    aModel.SetHash(a.Hash_id)
    aModel.SetPrice(a.Price)
    aModel.SetName(a.Name)
    lModel.SetLatitude(a.Gps.Lat)
    lModel.SetLongitude(a.Gps.Lon)
    aModel.SetLink(a.getLink())
    aModel.SetLocation(lModel)
    return aModel
}

func (a *advert) getLink() (string) {
    urlParts := &urlPartsDefinition {
        map[string]string{"1":"byt", "2":"dom", "666": "projekt", "3": "pozemka", "4": "komercni", "5": "ostatni"},
        map[string]string{"2":"1%2Bkk", "3": "1%2B1", "4": "2%2Bkk", "5": "2%2B1", "6": "3%2Bkk", "7": "3%2B1", "8": "4%2Bkk", "9": "4%2B1", "10":"5%2Bkk", "11":"5%2B1", "12": "6-a-vice", "16": "atypicky"},
        map[string]string{"1":"prodej", "2":"pronajem", "3": "drazby"} }
    url := domain
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_type_cb, 10)] + "/"
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_main_cb, 10)] + "/"
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_type_cb, 10)] + "/"
    url+= a.Seo.Locality + "/" + strconv.FormatInt(a.Hash_id, 10)
    return url;
}

