package email


import(
    "strconv"
    "net/mail"

    "../models"
)

func (ed *EmailDefinition) setFrom(email, name string) {
    ed.from = &mail.Address{name, email}
}

func (ed *EmailDefinition) setTo(email, name string) {
    ed.to = &mail.Address{name, email}
}

func (ed *EmailDefinition) setCC(email, name string) {
    //not implemented
}

func (ed *EmailDefinition) setSubject(subject string) {
    ed.subject = subject;
}

func (ed *EmailDefinition) setBody(body string) {
    ed.body = body;
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
    return "FlatSearch Agent: Found new advert " + strconv.Itoa(int(model.GetHash()));
}

func getBody(model *models.Advert) string {
    return "Hash: " + strconv.Itoa(int(model.GetHash())) + 
        "\nName: " + model.GetName() + 
        "\nLocality: " + model.GetLocality() + 
        "\nPrice: " + strconv.FormatFloat(model.GetPrice(), 'f', 2, 64) + 
        "\nUrl: " + model.GetLink() + 
        "\n";
}

