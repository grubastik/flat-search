package sreality

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type key string

// Internal map for converting data to detail url
type urlPartsDefinition struct {
	CategoryMainCb map[string]string
	CategorySubCb  map[string]string
	CategoryTypeCb map[string]string
}

// Detail defines top-level content of the response
type Detail struct {
	Properties []Property    `json:"items"`
	Locality   Property      `json:"locality"`
	Meta       string        `json:"meta_description"`
	Name       Property      `json:"name"`
	POI        []POILocation `json:"poi"`
	Price      Property      `json:"price_czk"`
	Text       Property      `json:"text"`
	Extra      Extra         `json:"_embedded"`
}

// Property defines property of the flat
type Property struct {
	Currency string      `json:"currency"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Unit     string      `json:"unit"`
	Value    interface{} `json:"value"`
}

// POILocation defines struct for storing nearby objects
type POILocation struct {
	Description  string  `json:"description"`
	Distance     float64 `json:"distance"`
	ImageURL     string  `json:"imgUrl"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Name         string  `json:"name"`
	WalkDistance int64   `json:"walkDistance"`
	Time         int64   `json:"time"`
}

// Extra stores extra info
type Extra struct {
	Images []Image `json:"images"`
	Seller Seller  `json:"seller"`
}

// Image stores info about image
type Image struct {
	Links Links `json:"_links"`
}

// Links stores info about image links
type Links struct {
	Main ImageURLDetail `json:"self"`
}

// ImageURLDetail stores info about image url
type ImageURLDetail struct {
	URL string `json:"href"`
}

// Seller stores info about realtor
type Seller struct {
	Email  string      `json:"email"`
	Image  string      `json:"image"`
	Phones []Phone     `json:"phones"`
	Name   string      `json:"user_name"`
	Extra  ExtraSeller `json:"_embedded"`
}

// Phone stores info about phone number
type Phone struct {
	Code   string `json:"code"`
	Number string `json:"number"`
	Type   string `json:"type"`
}

// ExtraSeller stores extra info about seller
type ExtraSeller struct {
	Company Company `json:"premise"`
}

// Company stores info about company where seller works
type Company struct {
	Address     string   `json:"address"`
	Description string   `json:"description"`
	Email       string   `json:"email"`
	ICO         int64    `json:"ico"`
	Location    Location `json:"locality"`
	Name        string   `json:"name"`
	Phones      []Phone  `json:"phones"`
	URL         string   `json:"www"`
}

// GetDetail sends request to the server to get details about advert
func GetDetail(ctx context.Context, a Advert) (*Detail, error) {
	rs := apiURL + "/" + strconv.FormatInt(a.HashID, 10)
	ctx = context.WithValue(ctx, requestKey, rs)
	request, err := http.NewRequest(http.MethodGet, rs, nil)
	if err != nil {
		return nil, err
	}

	// You can context separately for each request
	request = request.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(request)
	// You will get an error "net/http: request canceled" when request timeout exceeds limits
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, responseKey, body)

	var d = new(Detail)
	err = json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetDetailLink retrieves url to the details of the advert
func GetDetailLink(a Advert) string {
	//if detail url not set - return empty string
	if detailURL == "" {
		return ""
	}
	//define mapping betveen seo and url parameters
	urlParts := &urlPartsDefinition{
		map[string]string{"1": "byt", "2": "dom", "666": "projekt", "3": "pozemka", "4": "komercni", "5": "ostatni"},
		map[string]string{"2": "1%2Bkk", "3": "1%2B1", "4": "2%2Bkk", "5": "2%2B1", "6": "3%2Bkk", "7": "3%2B1", "8": "4%2Bkk", "9": "4%2B1", "10": "5%2Bkk", "11": "5%2B1", "12": "6-a-vice", "16": "atypicky"},
		map[string]string{"1": "prodej", "2": "pronajem", "3": "drazby"}}
	//build url
	url := detailURL
	url += urlParts.CategoryTypeCb[strconv.FormatInt(a.Seo.CategoryTypeCb, 10)] + "/"
	url += urlParts.CategoryMainCb[strconv.FormatInt(a.Seo.CategoryMainCb, 10)] + "/"
	url += urlParts.CategorySubCb[strconv.FormatInt(a.Seo.CategorySubCb, 10)] + "/"
	url += a.Seo.Locality + "/" + strconv.FormatInt(a.HashID, 10)
	return url
}
