package aggregator

import (
    "context"
    "time"
    "strconv"
    "reflect"

    "github.com/grubastik/flat-search/config"
    "github.com/grubastik/flat-search/models"
    "github.com/grubastik/flat-search/sreality"
)

type itos struct{
    Name  string
    Value string
}

// ProcessAdverts requests adverts from the service, populate advert structures, store them in DB and send email about new advert
func ProcessAdverts(c *config.Sreality) error {
    var l []*models.Advert;
    var cancel context.CancelFunc

    up := sreality.NewSrealityList(c.URL, c.URLDetail)
    PopulateSrealityParams(c, up)

    ctx := context.Background();
    ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    adverts, err := up.GetList(ctx)
    if err != nil {
        return err
    }

    for _, a := range *(adverts.GetAdverts()) {
        d, err := sreality.GetDetail(ctx, a)
        if (err != nil) {
            return err
        }
        l = append(l, ConvertSrealityToModel(a, *d))
    }
    err = saveToDb(l)
    if (err != nil) {
        return err
    }
    return snedEmails(l)
}

// PopulateSrealityParams populates info from config to sreality structure
func PopulateSrealityParams(c *config.Sreality, up *sreality.URLParamsList) {
    if c.RealityType > 0 {
        up.CategoryMainCb = int(c.RealityType)
    }
    if c.OperationType > 0 {
        up.CategoryTypeCb = int(c.OperationType)
    }
    if c.RealityOptions != nil {
        up.CategorySubCb = c.RealityOptions
    }
    if c.Country > 0 {
        up.LocalityCountryID = c.Country
    }
    if c.Region != nil {
        up.LocalityRegionID = c.Region
    }
    if c.District != nil {
        up.LocalityDistrictID = c.District
    }
    if c.PageResults > 0 {
        up.PerPage = c.PageResults
    }
    if c.EstateAge > 0 {
        up.EstateAge = c.EstateAge
    }
    if c.Square != nil {
        up.SetUsableArea(c.Square.Min, c.Square.Max)
    }
    if c.Price != nil {
        up.SetPriceRange(c.Price.Min, c.Price.Max)
    }
    if c.Furnished != nil {
        up.Furnished = c.Furnished
    }
}

// ConvertSrealityToModel builds model based on the data from response
func ConvertSrealityToModel(a sreality.Advert, d sreality.Detail) *models.Advert {
    var (
        am = models.NewAdvert()
        lm = new(models.Location)
        rm = new(models.Realtor)
        pm *models.Property
        im *models.Image
    )
    am.Locality = a.Locality
    am.HashID = a.HashID
    am.Price = a.Price
    am.Name = a.Name
    am.Link = sreality.GetDetailLink(a)
    am.Description = ""
    if (d.Text.Value != nil) {
        am.Description = d.Text.Value.(string)
    }
    lm.Lat = a.Gps.Lat
    lm.Lon = a.Gps.Lon
    am.Location = lm
    if &(d.Extra) != nil && &(d.Extra.Seller) != nil {
        rm.Name = d.Extra.Seller.Name
        rm.Email = d.Extra.Seller.Email
        if d.Extra.Seller.Phones != nil && len(d.Extra.Seller.Phones) > 0 {
            rm.Phone = d.Extra.Seller.Phones[0].Code + " " + d.Extra.Seller.Phones[0].Number
        }
        if &(d.Extra.Seller.Extra) != nil && &(d.Extra.Seller.Extra.Company) != nil {
            rm.Company = d.Extra.Seller.Extra.Company.Name
            if &(d.Extra.Seller.Extra.Company.Phones) != nil && len(d.Extra.Seller.Extra.Company.Phones) > 0 {
                rm.CompanyPhone = d.Extra.Seller.Extra.Company.Phones[0].Code + " " + d.Extra.Seller.Extra.Company.Phones[0].Number
            }
            rm.CompanyICO = strconv.FormatInt(d.Extra.Seller.Extra.Company.ICO, 10)
        }
        am.Realtor = *rm
    }
    if &(d.Extra) != nil && &(d.Extra.Images) != nil {
        for _, i := range d.Extra.Images {
            if &i != nil && &(i.Links) != nil && &(i.Links.Main) != nil && i.Links.Main.URL != "" {
                im = new(models.Image);
                im.URL = i.Links.Main.URL
                am.Images = append(am.Images, *im)
            }
        }
    }

    if &(d.Properties) != nil {
        for _, p := range d.Properties {
            pm = new(models.Property)
            pm.Name = p.Name
            switch p.Type {
                case "price_czk":
                if p.Value == nil {
                    pm.Value = ""
                } else {
                    pm.Value = p.Value.(string) + " " + p.Currency + " " + p.Unit
                }
                case "area":
                if p.Value == nil {
                    pm.Value = ""
                } else {
                    pm.Value = p.Value.(string) + p.Unit
                }
                case "string", "edited", "energy_efficiency_rating", "date":
                if p.Value == nil {
                    pm.Value = ""
                } else {
                    pm.Value = p.Value.(string)
                }
                case "boolean":
                pm.Value = "Ne"
                if p.Value != nil && p.Value.(bool) {
                    pm.Value = "Ano"
                }
                case "set":
                if p.Value == nil {
                    pm.Value = ""
                } else {
                    s := reflect.ValueOf(p.Value)
                    for i := 0; i < s.Len(); i++ {
                        is := s.Index(i).Interface().(map[string]interface{})
                        pm.Value += is["value"].(string) + "; "
                    }
                }
            }
            am.Properties = append(am.Properties, *pm)
        }
    }

    return am
}
