package sreality

import(
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/grubastik/flat-search/config"
    "github.com/grubastik/flat-search/email"
)

const configName = "sreality"

type rangeDefinition struct {
    min int
    max int
}

type UrlParams struct {
    url string
    category_main_cb int
    category_sub_cb *[]int
    category_type_cb int
    czk_price_summary_order2 *rangeDefinition
    estate_age int
    locality_district_id *[]int
    locality_region_id *[]int
    locality_country_id int
    per_page int
    usable_area *rangeDefinition
}

var conf *config.Config

func NewSreality(config *config.Config) *UrlParams {
    conf = config
    var urlParameters *UrlParams = new(UrlParams);
    moduleConfig := config.GetSreality();

    if (moduleConfig.RealityType > 0) {
        urlParameters.category_main_cb = int(moduleConfig.RealityType)
    }
    if (moduleConfig.OperationType > 0) {
        urlParameters.category_type_cb = int(moduleConfig.OperationType)
    }
    if (moduleConfig.RealityOptions != nil) {
        urlParameters.category_sub_cb = &(moduleConfig.RealityOptions)
    }
    if (moduleConfig.Country > 0) {
        urlParameters.locality_country_id = moduleConfig.Country
    }
    if (moduleConfig.Region != nil) {
        urlParameters.locality_region_id = &(moduleConfig.Region)
    }
    if (moduleConfig.District != nil) {
        urlParameters.locality_district_id = &(moduleConfig.District)
    }
    if (moduleConfig.PageResults > 0) {
        urlParameters.per_page = moduleConfig.PageResults
    }
    if (moduleConfig.EstateAge > 0) {
        urlParameters.estate_age = moduleConfig.EstateAge
    }
    if (moduleConfig.Square != nil) {
        urlParameters.SetUsableArea(moduleConfig.Square.Min, moduleConfig.Square.Max)
    }
    if (moduleConfig.Price != nil) {
        urlParameters.SetPriceRange(moduleConfig.Price.Min, moduleConfig.Price.Max)
    }
    if (moduleConfig.Url != "") {
        urlParameters.url = moduleConfig.Url
    }

    return urlParameters
}

func (up *UrlParams) SetUsableArea(min, max int) {
    up.usable_area = new(rangeDefinition);
    up.usable_area.min = min
    up.usable_area.max = max
}

func (up *UrlParams) SetPriceRange(min, max int) {
    up.czk_price_summary_order2 = new(rangeDefinition);
    up.czk_price_summary_order2.min = min
    up.czk_price_summary_order2.max = max
}

func (up *UrlParams) GetRequest() (request string) {
    requestParts := []string{}
    if up.category_main_cb != 0 {
        requestParts = append(requestParts, "category_main_cb=" + strconv.Itoa(up.category_main_cb))
    }
    if up.category_type_cb != 0 {
        requestParts = append(requestParts, "category_type_cb=" + strconv.Itoa(up.category_type_cb))
    }

    if len(*up.category_sub_cb) > 0 {
        requestParts = append(requestParts, "category_sub_cb=" + getStringFromIntSlice(up.category_sub_cb))
    }
    if up.locality_country_id != 0 {
        requestParts = append(requestParts, "locality_country_id=" + strconv.Itoa(up.locality_country_id))
    }
    if len(*up.locality_region_id) > 0 {
        requestParts = append(requestParts, "locality_region_id=" + getStringFromIntSlice(up.locality_region_id))
    }
    if len(*up.locality_district_id) > 0 {
        requestParts = append(requestParts, "locality_district_id=" + getStringFromIntSlice(up.locality_district_id))
    }
    if up.per_page != 0 {
        requestParts = append(requestParts, "per_page=" + strconv.Itoa(up.per_page))
    }
    if up.estate_age != 0 {
        requestParts = append(requestParts, "estate_age=" + strconv.Itoa(up.estate_age))
    }
    if up.usable_area != nil {
        requestParts = append(requestParts, "usable_area=" + getStringFromRange(up.usable_area))
    }
    if up.czk_price_summary_order2 != nil {
        requestParts = append(requestParts, "czk_price_summary_order2=" + getStringFromRange(up.czk_price_summary_order2))
    }
    request = ""
    if len(requestParts) > 0 {
        request += "?" + strings.Join(requestParts, "&")
    }
    return up.url + request
}

func getStringFromIntSlice(val *[]int) string {
    request := "";
    for i := 0; i < len(*val); i++ {
        request+= strconv.Itoa((*val)[i]) + "|"
    }
    request = strings.TrimRight(request, "|")
    return request
}

func getStringFromRange(val *rangeDefinition) string {
    return strconv.Itoa(val.min) + "|" + strconv.Itoa(val.max)
}

func (up *UrlParams) MakeRequest() (*Body, error) {
    resp, err := http.Get(up.GetRequest())
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var adverts = new(Body)
    err = json.Unmarshal(body, adverts)
    if err != nil {
        return nil, err
    }
    return adverts, nil;
}

func (up *UrlParams) ProcessAdverts() error {
    adverts, err := up.MakeRequest()
    if err != nil {
        return err
    }
    for _, advert := range *(adverts.GetAdverts()) {
        aModel := advert.ConvertToModel()
        exists, err := aModel.ExistsInDbByHashId()
        if err != nil {
            return err
        }
        if (!exists) {
            //record to db
            err = aModel.Insert()
            if err != nil {
                return err
            }
            //sendemail
            err = email.Conn.Send(email.NewEmail(aModel))
            if err != nil {
                return err
            }
        }
    }
    return nil;
}
