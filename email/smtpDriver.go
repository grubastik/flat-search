package email

import (
    "fmt"
    "net"
    "net/mail"
	"net/smtp"
    "crypto/tls"
    "strconv"
    "./../error"
    "./../config"
)

const configName = "email"

type EmailDefinition struct {
    from *mail.Address
    to *mail.Address
    subject string
    body string
    headers map[string]string
    connection *smtp.Client
    server string
    port int
    tlsEnabled bool
    username string
    password string
}

func New(config *config.Config) (*EmailDefinition) {
    var emailDefinition *EmailDefinition = new(EmailDefinition);
    moduleConfig := config.GetEmail();

    if (len(moduleConfig.To) > 0) {
        emailDefinition.setTo(moduleConfig.To, "");
    }
    if (len(moduleConfig.From) > 0) {
        emailDefinition.setFrom(moduleConfig.From, "");
    }
    if (len(moduleConfig.Server) > 0) {
        emailDefinition.server = moduleConfig.Server;
    }
    if (moduleConfig.TlsPort > 0) {
        emailDefinition.port = moduleConfig.TlsPort;
    }
    if (moduleConfig.Tls) {
        emailDefinition.tlsEnabled = moduleConfig.Tls;
    }
    if (len(moduleConfig.Username) > 0) {
        emailDefinition.username = moduleConfig.Username;
    }
    if (len(moduleConfig.Password) > 0) {
        emailDefinition.password = moduleConfig.Password;
    }
    return emailDefinition;
}

func (ed *EmailDefinition) getMessage() (string) {
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

func (ed *EmailDefinition) Send() {
    // Setup message
    message := ed.getMessage();

    ed.makeSmtpClient()
    defer ed.connection.Quit()
    ed.authenticate()

    // To && From
    err := ed.connection.Mail(ed.from.Address);
    error.DebugError(err)
    err = ed.connection.Rcpt(ed.to.Address);
    error.DebugError(err)

    // Data
    w, err := ed.connection.Data()
    error.DebugError(err)

    _, err = w.Write([]byte(message))
    error.DebugError(err)

    err = w.Close()
    error.DebugError(err)
}

func (ed *EmailDefinition) authenticate() {
    host := GetHost(ed.server + ":" + strconv.Itoa(ed.port))
    auth := smtp.PlainAuth("", ed.username, ed.password, host)
    err := ed.connection.Auth(auth);
    error.DebugError(err)
}

func GetHost(servername string) string {
    host, _, _ := net.SplitHostPort(servername)
    return host
}

func (ed *EmailDefinition) makeSmtpClient() {
    host := GetHost(ed.server + ":" + strconv.Itoa(ed.port))

    // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }
    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", ed.server + ":" + strconv.Itoa(ed.port), tlsconfig)
    error.DebugError(err)

    c, err := smtp.NewClient(conn, host)
    error.DebugError(err)
    
    ed.connection = c
}
