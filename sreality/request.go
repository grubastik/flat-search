package sreality

import(
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "./../error"
    "./../config"
)

const configName = "sreality"

type rangeDefinition struct {
    min int
    max int
}

type UrlParams struct {
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

func New(config *config.Config) (*UrlParams) {
    var urlParameters *UrlParams = new(UrlParams);
    moduleConfig := config.GetSreality();

    if (moduleConfig.RealityType > 0) {
        urlParameters.SetRealityType(int(moduleConfig.RealityType))
    }
    if (moduleConfig.OperationType > 0) {
        urlParameters.SetOperationType(int(moduleConfig.OperationType))
    }
    if (moduleConfig.RealityOptions != nil) {
        urlParameters.SetRealityOptions(&(moduleConfig.RealityOptions))
    }
    if (moduleConfig.Country > 0) {
        urlParameters.SetCountry(moduleConfig.Country)
    }
    if (moduleConfig.Region != nil) {
        urlParameters.SetRegion(&(moduleConfig.Region))
    }
    if (moduleConfig.District != nil) {
        urlParameters.SetDistrict(&(moduleConfig.District))
    }
    if (moduleConfig.PageResults > 0) {
        urlParameters.SetPagedResults(moduleConfig.PageResults)
    }
    if (moduleConfig.EstateAge > 0) {
        urlParameters.SetEstateAge(moduleConfig.PageResults)
    }
    if (moduleConfig.Square != nil) {
        urlParameters.SetUsableArea(moduleConfig.Square.Min, moduleConfig.Square.Max)
    }
    if (moduleConfig.Price != nil) {
        urlParameters.SetPriceRange(moduleConfig.Price.Min, moduleConfig.Price.Max)
    }

    return urlParameters
}

func (up *UrlParams) SetRealityType(realityType int) {
    up.category_main_cb = realityType
}

func (up *UrlParams) SetOperationType(operationType int) {
    up.category_type_cb = operationType
}

func (up *UrlParams) SetRealityOptions(realityOptions *[]int) {
    up.category_sub_cb = realityOptions
}

func (up *UrlParams) SetCountry(country int) {
    up.locality_country_id = country
}

func (up *UrlParams) SetRegion(region *[]int) {
    up.locality_region_id = region//112
}

func (up *UrlParams) SetDistrict(districts *[]int) {
    up.locality_district_id = districts//112
}

func (up *UrlParams) SetPagedResults(resultsPerPage int) {
    up.per_page = resultsPerPage
}

func (up *UrlParams) SetEstateAge(age int) {
    up.estate_age = age
}

func (up *UrlParams) SetUsableArea(min, max int) {
    up.usable_area = getRange(min, max)
}

func (up *UrlParams) SetPriceRange(min, max int) {
    up.czk_price_summary_order2 = getRange(min, max)
}

func (up *UrlParams) GetRequest() (request string) {
    requestParts := []string{}
    if up.category_main_cb != 0 {
        requestParts = append(requestParts, "category_main_cb=" + getStringFromInt(up.category_main_cb))
    }
    if up.category_type_cb != 0 {
        requestParts = append(requestParts, "category_type_cb=" + getStringFromInt(up.category_type_cb))
    }

    if len(*up.category_sub_cb) > 0 {
        requestParts = append(requestParts, "category_sub_cb=" + getStringFromIntSlice(up.category_sub_cb))
    }
    if up.locality_country_id != 0 {
        requestParts = append(requestParts, "locality_country_id=" + getStringFromInt(up.locality_country_id))
    }
    if len(*up.locality_region_id) > 0 {
        requestParts = append(requestParts, "locality_region_id=" + getStringFromIntSlice(up.locality_region_id))
    }
    if len(*up.locality_district_id) > 0 {
        requestParts = append(requestParts, "locality_district_id=" + getStringFromIntSlice(up.locality_district_id))
    }
    if up.per_page != 0 {
        requestParts = append(requestParts, "per_page=" + getStringFromInt(up.per_page))
    }
    if up.estate_age != 0 {
        requestParts = append(requestParts, "estate_age=" + getStringFromInt(up.estate_age))
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
    return request
}

func getRange(min, max int) *rangeDefinition {
    rangeMy := new(rangeDefinition);
    rangeMy.min = min
    rangeMy.max = max
    return rangeMy
}

func getStringFromInt(val int) string {
    return strconv.Itoa(val)
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

func (up *UrlParams) MakeRequest() (*Body) {
    resp, err := http.Get("https://www.sreality.cz/api/cs/v2/estates" + up.GetRequest())
    error.DebugError(err)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    error.DebugError(err)

    var adverts = new(Body)
    err = json.Unmarshal(body, adverts)
    error.DebugError(err)
    return adverts;
}
