package email

import (
	"crypto/tls"
	"github.com/grubastik/flat-search/config"
	"net"
	"net/smtp"
	"strconv"
)

const configName = "email"

// Connection defines information about connection opened to smtp server
type Connection struct {
	connection *smtp.Client
	server     string
	port       int
	tlsEnabled bool
	username   string
	password   string
}

var emailConf *config.Email

// Conn contains information about connection opened to smtp server
var Conn *Connection

// NewConnection populates struct Connection and opens connection to the smtp server
func NewConnection(config *config.Config) *Connection {
	Conn = new(Connection)
	emailConf = config.GetEmail()

	if emailConf != nil && emailConf.Server != "" {
		Conn.server = emailConf.Server
	}
	if emailConf != nil && emailConf.TLSPort > 0 {
		Conn.port = emailConf.TLSPort
	}
	if emailConf != nil && emailConf.TLS {
		Conn.tlsEnabled = emailConf.TLS
	}
	if emailConf != nil && emailConf.Username != "" {
		Conn.username = emailConf.Username
	}
	if emailConf != nil && emailConf.Password != "" {
		Conn.password = emailConf.Password
	}
	return Conn
}

// GetHost accepts string with service address and get host name from it
func GetHost(servername string) (string, error) {
	host, _, err := net.SplitHostPort(servername)
	return host, err
}

// Send send email through the opened connection
func (ec *Connection) Send(ed *Definition) error {
	// Setup message
	message := ed.getMessage()

	ec.makeSMTPClient()
	defer ec.connection.Quit()
	ec.authenticate()

	// To && From
	err := ec.connection.Mail(ed.from.Address)
	if err != nil {
		return err
	}
	err = ec.connection.Rcpt(ed.to.Address)
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

func (ec *Connection) authenticate() error {
	host, err := GetHost(ec.server + ":" + strconv.Itoa(ec.port))
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", ec.username, ec.password, host)
	err = ec.connection.Auth(auth)
	if err != nil {
		return err
	}
	return nil
}

func (ec *Connection) makeSMTPClient() error {
	host, err := GetHost(ec.server + ":" + strconv.Itoa(ec.port))
	if err != nil {
		return err
	}

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", ec.server+":"+strconv.Itoa(ec.port), tlsconfig)
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
