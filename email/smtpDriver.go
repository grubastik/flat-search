package email

import (
    "fmt"
    "net"
    "net/mail"
	"net/smtp"
    "crypto/tls"
    "strconv"
    "github.com/grubastik/flat-search/config"
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

func NewEmail(config *config.Config) (*EmailDefinition) {
    var emailDefinition *EmailDefinition = new(EmailDefinition);
    moduleConfig := config.GetEmail();

    if len(moduleConfig.To) > 0 {
        emailDefinition.to = &mail.Address{"", moduleConfig.To}
    }
    if len(moduleConfig.From) > 0 {
        emailDefinition.from = &mail.Address{"", moduleConfig.From}
    }
    if len(moduleConfig.Server) > 0 {
        emailDefinition.server = moduleConfig.Server;
    }
    if moduleConfig.TlsPort > 0 {
        emailDefinition.port = moduleConfig.TlsPort;
    }
    if moduleConfig.Tls {
        emailDefinition.tlsEnabled = moduleConfig.Tls;
    }
    if len(moduleConfig.Username) > 0 {
        emailDefinition.username = moduleConfig.Username;
    }
    if len(moduleConfig.Password) > 0 {
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

func (ed *EmailDefinition) Send() (error) {
    // Setup message
    message := ed.getMessage();

    ed.makeSmtpClient()
    defer ed.connection.Quit()
    ed.authenticate()

    // To && From
    err := ed.connection.Mail(ed.from.Address);
    if err != nil {
        return err
    }
    err = ed.connection.Rcpt(ed.to.Address);
    if err != nil {
        return err
    }

    // Data
    w, err := ed.connection.Data()
    if err != nil {
        return err
    }

    _, err = w.Write([]byte(message))
    if err != nil {
        return err
    }

    err = w.Close()
    if err != nil {
        return err
    }
    return nil
}

func (ed *EmailDefinition) authenticate() (error) {
    host, err := GetHost(ed.server + ":" + strconv.Itoa(ed.port))
    if err != nil {
        return err
    }
    auth := smtp.PlainAuth("", ed.username, ed.password, host)
    err = ed.connection.Auth(auth);
    if err != nil {
        return err
    }
    return nil
}

func GetHost(servername string) (string, error) {
    host, _, err := net.SplitHostPort(servername)
    return host, err
}

func (ed *EmailDefinition) makeSmtpClient() (error) {
    host, err := GetHost(ed.server + ":" + strconv.Itoa(ed.port))
    if err != nil {
        return err
    }

    // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }
    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", ed.server + ":" + strconv.Itoa(ed.port), tlsconfig)
    if err != nil {
        return err
    }

    c, err := smtp.NewClient(conn, host)
    if err != nil {
        return err
    }
    
    ed.connection = c
    return nil
}
