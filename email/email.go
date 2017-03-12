package email


import(
    "strconv"

    "github.com/grubastik/flat-search/models"
)

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

