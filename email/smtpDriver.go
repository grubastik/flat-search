package email

import (
    "net"
	"net/smtp"
    "crypto/tls"
    "strconv"
    "github.com/grubastik/flat-search/config"
)

const configName = "email"

type EmailConnection struct {
    connection *smtp.Client
    server string
    port int
    tlsEnabled bool
    username string
    password string
}

var emailConf *config.Email

var Conn *EmailConnection

func NewEmailConnection(config *config.Config) *EmailConnection {
    var ec *EmailConnection = new(EmailConnection);
    emailConf = config.GetEmail();

    if len(emailConf.Server) > 0 {
        ec.server = emailConf.Server;
    }
    if emailConf.TlsPort > 0 {
        ec.port = emailConf.TlsPort;
    }
    if emailConf.Tls {
        ec.tlsEnabled = emailConf.Tls;
    }
    if len(emailConf.Username) > 0 {
        ec.username = emailConf.Username;
    }
    if len(emailConf.Password) > 0 {
        ec.password = emailConf.Password;
    }
    return ec;
}

func (ec *EmailConnection) Send(ed *EmailDefinition) error {
    // Setup message
    message := ed.getMessage();

    ec.makeSmtpClient()
    defer ec.connection.Quit()
    ec.authenticate()

    // To && From
    err := ec.connection.Mail(ed.from.Address);
    if err != nil {
        return err
    }
    err = ec.connection.Rcpt(ed.to.Address);
    if err != nil {
        return err
    }

    // Data
    w, err := ec.connection.Data()
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

func (ec *EmailConnection) authenticate() error {
    host, err := GetHost(ec.server + ":" + strconv.Itoa(ec.port))
    if err != nil {
        return err
    }
    auth := smtp.PlainAuth("", ec.username, ec.password, host)
    err = ec.connection.Auth(auth);
    if err != nil {
        return err
    }
    return nil
}

func GetHost(servername string) (string, error) {
    host, _, err := net.SplitHostPort(servername)
    return host, err
}

func (ec *EmailConnection) makeSmtpClient() error {
    host, err := GetHost(ec.server + ":" + strconv.Itoa(ec.port))
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
    conn, err := tls.Dial("tcp", ec.server + ":" + strconv.Itoa(ec.port), tlsconfig)
    if err != nil {
        return err
    }

    c, err := smtp.NewClient(conn, host)
    if err != nil {
        return err
    }
    
    ec.connection = c
    return nil
}
