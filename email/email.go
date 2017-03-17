package email


import(
    "fmt"
    "strconv"
    "net/mail"
    "github.com/grubastik/flat-search/models"
)

type EmailDefinition struct {
    from *mail.Address
    to *mail.Address
    subject string
    body string
    headers map[string]string
}

func NewEmail(model *models.Advert) *EmailDefinition {
    var ed *EmailDefinition = new(EmailDefinition);
    if len(emailConf.To) > 0 {
        ed.to = &mail.Address{"", emailConf.To}
    }
    if len(emailConf.From) > 0 {
        ed.from = &mail.Address{"", emailConf.From}
    }
    ed.PrepareData(model)
    return ed
}

func (ed *EmailDefinition) getMessage() string {
    message := ""
    for k,v := range ed.headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += fmt.Sprintf("%s: %s\r\n", "From", ed.from.String())
    message += fmt.Sprintf("%s: %s\r\n", "To", ed.to.String())
    message += fmt.Sprintf("%s: %s\r\n", "Subject", ed.subject)
    message += "\r\n" + ed.body
    return message
}

func (ed *EmailDefinition) clearHeaders() {
    ed.headers = make(map[string]string)
}

func (ed *EmailDefinition) addHeader(name, value string) {
    ed.headers[name] = value;
}

func (ed *EmailDefinition) PrepareData(model *models.Advert) {
    ed.subject = getSubject(model);
    ed.body = getBody(model);
}

func getSubject(model *models.Advert) string {
    return "FlatSearch Agent: Found new advert " + strconv.Itoa(int(model.HashId));
}

func getBody(model *models.Advert) string {
    return "Hash: " + strconv.Itoa(int(model.HashId)) + 
        "\nName: " + model.Name + 
        "\nLocality: " + model.Locality + 
        "\nPrice: " + strconv.FormatFloat(model.Price, 'f', 2, 64) + 
        "\nUrl: " + model.Link + 
        "\n";
}

