package sreality

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "context"
)

const (
    requestKey key = "request"
    responseKey key = "response"
)

var (
    detailURL string
    apiURL string
)

type rangeDefinition struct {
    min int
    max int
}

// URLParamsList defines all vars to build list url
type URLParamsList struct {
    CategoryMainCb        int
    CategorySubCb         []int
    CategoryTypeCb        int
    EstateAge             int
    LocalityDistrictID    []int
    LocalityRegionID      []int
    LocalityCountryID     int
    PerPage               int
    usableArea            *rangeDefinition
    czkPriceSummaryOrder2 *rangeDefinition
}

//Location stores geographical coordinates
type Location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

//Seo stores info for building detail url
type Seo struct {
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
    Seo      *Seo      `json:"seo"`
    Gps      *Location `json:"gps"`
}

// Adverts contains list of adverts
type Adverts struct {
    Estates *[]Advert `json:"estates"`
}

// Body defines top-level content of the response
type Body struct {
    Embedded *Adverts `json:"_embedded"`
}

// SetUsableArea sets area range of the flat
func (up *URLParamsList) SetUsableArea(min, max int) {
    up.usableArea = new(rangeDefinition)
    up.usableArea.min = min
    up.usableArea.max = max
}

// SetPriceRange sets price range of the flat
func (up *URLParamsList) SetPriceRange(min, max int) {
    up.czkPriceSummaryOrder2 = new(rangeDefinition)
    up.czkPriceSummaryOrder2.min = min
    up.czkPriceSummaryOrder2.max = max
}

// NewSrealityList creates base structure and init vars of the package
func NewSrealityList(au, du string) *URLParamsList {
    up := new(URLParamsList)
    detailURL = du
    apiURL = au
    return up
}

// GetRequest build url to get list of the adverts
func (up *URLParamsList) GetRequest() (request string) {
    requestParts := []string{}
    if up.CategoryMainCb != 0 {
        requestParts = append(requestParts, "category_main_cb="+strconv.Itoa(up.CategoryMainCb))
    }
    if up.CategoryTypeCb != 0 {
        requestParts = append(requestParts, "category_type_cb="+strconv.Itoa(up.CategoryTypeCb))
    }

    if up.CategorySubCb != nil && len(up.CategorySubCb) > 0 {
        requestParts = append(requestParts, "category_sub_cb="+getStringFromIntSlice(up.CategorySubCb))
    }
    if up.LocalityCountryID != 0 {
        requestParts = append(requestParts, "locality_country_id="+strconv.Itoa(up.LocalityCountryID))
    }
    if up.LocalityRegionID != nil && len(up.LocalityRegionID) > 0 {
        requestParts = append(requestParts, "locality_region_id="+getStringFromIntSlice(up.LocalityRegionID))
    }
    if up.LocalityDistrictID != nil && len(up.LocalityDistrictID) > 0 {
        requestParts = append(requestParts, "locality_district_id="+getStringFromIntSlice(up.LocalityDistrictID))
    }
    if up.PerPage != 0 {
        requestParts = append(requestParts, "per_page="+strconv.Itoa(up.PerPage))
    }
    if up.EstateAge != 0 {
        requestParts = append(requestParts, "estate_age="+strconv.Itoa(up.EstateAge))
    }
    if up.usableArea != nil {
        requestParts = append(requestParts, "usable_area="+getStringFromRange(up.usableArea))
    }
    if up.czkPriceSummaryOrder2 != nil {
        requestParts = append(requestParts, "czk_price_summary_order2="+getStringFromRange(up.czkPriceSummaryOrder2))
    }
    request = ""
    if len(requestParts) > 0 {
        request += "?" + strings.Join(requestParts, "&")
    }
    return apiURL + request
}

// GetList requestd list of the adverts from the service
func (up *URLParamsList) GetList(ctx context.Context) (*Body, error) {
    rs := up.GetRequest()
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

    var adverts = new(Body)
    err = json.Unmarshal(body, adverts)
    if err != nil {
        return nil, err
    }
    return adverts, nil
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

func getStringFromIntSlice(val []int) string {
    request := ""
    for i := 0; i < len(val); i++ {
        request += strconv.Itoa(val[i]) + "|"
    }
    request = strings.TrimRight(request, "|")
    return request
}

func getStringFromRange(val *rangeDefinition) string {
    return strconv.Itoa(val.min) + "|" + strconv.Itoa(val.max)
}

