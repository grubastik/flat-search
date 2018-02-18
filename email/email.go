package email

import (
	"fmt"
	"github.com/grubastik/flat-search/models"
	"net/mail"
	"strconv"
	"errors"
)

// Definition stores email info
type Definition struct {
	from    *mail.Address
	to      []*mail.Address
	subject string
	body    string
	headers map[string]string
}

// NewEmail creates new email struct and populate it with values stored in config
func NewEmail(model *models.Advert) (*Definition, error) {
	ed := new(Definition)
	if emailConf != nil && emailConf.To != "" {
		to, err := mail.ParseAddressList(emailConf.To)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to parse email list: %s", err.Error()))
		}
		ed.to = to
	}
	if emailConf != nil && emailConf.From != "" {
		ed.from = &mail.Address{"", emailConf.From}
	}
	ed.PrepareData(model)
	return ed, nil
}

// PrepareData populates struct with imail info based on the model data
func (ed *Definition) PrepareData(model *models.Advert) {
	ed.subject = getSubject(model)
	ed.body = getBody(model)
}

func (ed *Definition) getMessage() string {
	message := ""
	to := ""
	for k, v := range ed.headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	for _, t := range ed.to {
		to += fmt.Sprintf("%s", t.String())
		to += ";"
	}
	message += fmt.Sprintf("%s: %s\r\n", "From", ed.from.String())
	message += fmt.Sprintf("%s: %s\r\n", "To", to[:len(to) - 1])
	message += fmt.Sprintf("%s: %s\r\n", "Subject", ed.subject)
	message += "\r\n" + ed.body
	return message
}

func (ed *Definition) clearHeaders() {
	ed.headers = make(map[string]string)
}

func (ed *Definition) addHeader(name, value string) {
	ed.headers[name] = value
}

func getSubject(model *models.Advert) string {
	return "FlatSearch Agent: Found new advert " + strconv.Itoa(int(model.HashID))
}

func getBody(model *models.Advert) string {
	body := "Hash: " + strconv.Itoa(int(model.HashID)) +
		"\nName: " + model.Name +
		"\nLocality: " + model.Locality +
		"\nPrice: " + strconv.FormatFloat(model.Price, 'f', 2, 64) +
		"\nUrl: " + model.Link +
		"\n\n " + model.Description +
		"\n\nRealtor: " + model.Realtor.Name + "(phone: " + model.Realtor.Phone + "; email: " + model.Realtor.Email + ")"

	if (model.Properties != nil) {
	    body+= "\n\n Properties:"
	    for _,p := range model.Properties {
	        body+= "\n" + p.Name + ": " + p.Value
	    }
    }
    
	if (model.Images != nil) {
        body+= "\n\n Images:"
	    for k,i := range model.Images {
	        body+= "\n\n - " + strconv.FormatInt(int64(k + 1), 10) + " - " + i.URL
	    }
    }
    
    body+= "\n"
    return body
}
