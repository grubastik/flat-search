package sreality

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/grubastik/flat-search/config"
	"github.com/grubastik/flat-search/email"
)

const configName = "sreality"

type rangeDefinition struct {
	min int
	max int
}

// URLParams defines all vars to build list url
type URLParams struct {
	url                   string
	categoryMainCb        int
	categorySubCb         *[]int
	categoryTypeCb        int
	czkPriceSummaryOrder2 *rangeDefinition
	estateAge             int
	localityDistrictID    *[]int
	localityRegionID      *[]int
	localityCountryID     int
	perPage               int
	usableArea            *rangeDefinition
}

var conf *config.Config

// NewSreality creates new structure and populate it with values from config
func NewSreality(config *config.Config) *URLParams {
	conf = config
	urlParameters := new(URLParams)
	moduleConfig := config.GetSreality()

	if moduleConfig.RealityType > 0 {
		urlParameters.categoryMainCb = int(moduleConfig.RealityType)
	}
	if moduleConfig.OperationType > 0 {
		urlParameters.categoryTypeCb = int(moduleConfig.OperationType)
	}
	if moduleConfig.RealityOptions != nil {
		urlParameters.categorySubCb = &(moduleConfig.RealityOptions)
	}
	if moduleConfig.Country > 0 {
		urlParameters.localityCountryID = moduleConfig.Country
	}
	if moduleConfig.Region != nil {
		urlParameters.localityRegionID = &(moduleConfig.Region)
	}
	if moduleConfig.District != nil {
		urlParameters.localityDistrictID = &(moduleConfig.District)
	}
	if moduleConfig.PageResults > 0 {
		urlParameters.perPage = moduleConfig.PageResults
	}
	if moduleConfig.EstateAge > 0 {
		urlParameters.estateAge = moduleConfig.EstateAge
	}
	if moduleConfig.Square != nil {
		urlParameters.SetUsableArea(moduleConfig.Square.Min, moduleConfig.Square.Max)
	}
	if moduleConfig.Price != nil {
		urlParameters.SetPriceRange(moduleConfig.Price.Min, moduleConfig.Price.Max)
	}
	if moduleConfig.URL != "" {
		urlParameters.url = moduleConfig.URL
	}

	return urlParameters
}

// SetUsableArea sets area range of the flat
func (up *URLParams) SetUsableArea(min, max int) {
	up.usableArea = new(rangeDefinition)
	up.usableArea.min = min
	up.usableArea.max = max
}

// SetPriceRange sets price range of the flat
func (up *URLParams) SetPriceRange(min, max int) {
	up.czkPriceSummaryOrder2 = new(rangeDefinition)
	up.czkPriceSummaryOrder2.min = min
	up.czkPriceSummaryOrder2.max = max
}

// GetRequest build url to get list of the adverts
func (up *URLParams) GetRequest() (request string) {
	requestParts := []string{}
	if up.categoryMainCb != 0 {
		requestParts = append(requestParts, "category_main_cb="+strconv.Itoa(up.categoryMainCb))
	}
	if up.categoryTypeCb != 0 {
		requestParts = append(requestParts, "category_type_cb="+strconv.Itoa(up.categoryTypeCb))
	}

	if len(*up.categorySubCb) > 0 {
		requestParts = append(requestParts, "category_sub_cb="+getStringFromIntSlice(up.categorySubCb))
	}
	if up.localityCountryID != 0 {
		requestParts = append(requestParts, "locality_country_id="+strconv.Itoa(up.localityCountryID))
	}
	if len(*up.localityRegionID) > 0 {
		requestParts = append(requestParts, "locality_region_id="+getStringFromIntSlice(up.localityRegionID))
	}
	if len(*up.localityDistrictID) > 0 {
		requestParts = append(requestParts, "locality_district_id="+getStringFromIntSlice(up.localityDistrictID))
	}
	if up.perPage != 0 {
		requestParts = append(requestParts, "per_page="+strconv.Itoa(up.perPage))
	}
	if up.estateAge != 0 {
		requestParts = append(requestParts, "estate_age="+strconv.Itoa(up.estateAge))
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
	return up.url + request
}

// MakeRequest requestd list of the adverts from the service
func (up *URLParams) MakeRequest() (*Body, error) {
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
	return adverts, nil
}

// ProcessAdverts requests adverts from the service, populate advert structures, store them in DB and send email about new advert
func (up *URLParams) ProcessAdverts() error {
	adverts, err := up.MakeRequest()
	if err != nil {
		return err
	}
	for _, advert := range *(adverts.GetAdverts()) {
		aModel := advert.ConvertToModel()
		exists, err := aModel.ExistsInDbByHashID()
		if err != nil {
			return err
		}
		if !exists {
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
	return nil
}

func getStringFromIntSlice(val *[]int) string {
	request := ""
	for i := 0; i < len(*val); i++ {
		request += strconv.Itoa((*val)[i]) + "|"
	}
	request = strings.TrimRight(request, "|")
	return request
}

func getStringFromRange(val *rangeDefinition) string {
	return strconv.Itoa(val.min) + "|" + strconv.Itoa(val.max)
}
