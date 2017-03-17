
package sreality

import(
    "strconv"
    "github.com/grubastik/flat-search/models"
)

type urlPartsDefinition struct {
    Category_main_cb map[string]string
    Category_sub_cb map[string]string
    Category_type_cb map[string]string
}

type location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type seo struct {
    Category_main_cb int64 `json:"category_main_cb"`
    Category_sub_cb int64 `json:"category_sub_cb"`
    Category_type_cb int64 `json:"category_type_cb"`
    Locality string `json:"locality"`
}

type advert struct {
    Locality string `json:"locality"`
    Hash_id int64 `json:"hash_id"`
    Price float64 `json:"price"`
    Name string `json:"name"`
    Seo *seo `json:"seo"`
    Gps *location `json:"gps"`
}

type adverts struct {
    Estates *[]advert `json:"estates"`
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

func (a *advert) ConvertToModel() *models.Advert {
    var (
        aModel = models.NewAdvert()
        lModel = new(models.Location)
    )
    aModel.Locality = a.Locality
    aModel.HashId = a.Hash_id
    aModel.Price = a.Price
    aModel.Name = a.Name
    lModel.Lat = a.Gps.Lat
    lModel.Lon = a.Gps.Lon
    aModel.Link = a.getLink()
    aModel.Location = lModel
    return aModel
}

func (a *advert) getLink() string {
    urlParts := &urlPartsDefinition {
        map[string]string{"1":"byt", "2":"dom", "666": "projekt", "3": "pozemka", "4": "komercni", "5": "ostatni"},
        map[string]string{"2":"1%2Bkk", "3": "1%2B1", "4": "2%2Bkk", "5": "2%2B1", "6": "3%2Bkk", "7": "3%2B1", "8": "4%2Bkk", "9": "4%2B1", "10":"5%2Bkk", "11":"5%2B1", "12": "6-a-vice", "16": "atypicky"},
        map[string]string{"1":"prodej", "2":"pronajem", "3": "drazby"} }
    url := conf.GetSreality().UrlDetail
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_type_cb, 10)] + "/"
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_main_cb, 10)] + "/"
    url+= urlParts.Category_type_cb[strconv.FormatInt(a.Seo.Category_type_cb, 10)] + "/"
    url+= a.Seo.Locality + "/" + strconv.FormatInt(a.Hash_id, 10)
    return url;
}

