package email

import (
	"fmt"
	"github.com/grubastik/flat-search/models"
	"net/mail"
	"strconv"
)

// Definition stores email info
type Definition struct {
	from    *mail.Address
	to      *mail.Address
	subject string
	body    string
	headers map[string]string
}

// NewEmail creates new email struct and populate it with values stored in config
func NewEmail(model *models.Advert) *Definition {
	ed := new(Definition)
	if emailConf != nil && len(emailConf.To) > 0 {
		ed.to = &mail.Address{"", emailConf.To}
	}
	if emailConf != nil && len(emailConf.From) > 0 {
		ed.from = &mail.Address{"", emailConf.From}
	}
	ed.PrepareData(model)
	return ed
}

// PrepareData populates struct with imail info based on the model data
func (ed *Definition) PrepareData(model *models.Advert) {
	ed.subject = getSubject(model)
	ed.body = getBody(model)
}

func (ed *Definition) getMessage() string {
	message := ""
	for k, v := range ed.headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += fmt.Sprintf("%s: %s\r\n", "From", ed.from.String())
	message += fmt.Sprintf("%s: %s\r\n", "To", ed.to.String())
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
	return "Hash: " + strconv.Itoa(int(model.HashID)) +
		"\nName: " + model.Name +
		"\nLocality: " + model.Locality +
		"\nPrice: " + strconv.FormatFloat(model.Price, 'f', 2, 64) +
		"\nUrl: " + model.Link +
		"\n"
}
