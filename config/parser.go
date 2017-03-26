package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config describes config options
type Config struct {
	Db       *Db
	Email    *Email
	Sreality *Sreality
}

// Db describes config options for DB
type Db struct {
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Email describes config options for email sending
type Email struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Server   string `json:"server"`
	TLS      bool   `json:"tls"`
	TLSPort  int    `json:"tls_port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Sreality describes config options for the flat search request
type Sreality struct {
	URL            string    `json:"url"`
	URLDetail      string    `json:"url_detail"`
	RealityType    int       `json:"reality_type"`
	OperationType  int       `json:"operation_type"`
	RealityOptions []int     `json:"reality_options"`
	Country        int       `json:"country"`
	Region         []int     `json:"region"`
	District       []int     `json:"district"`
	PageResults    int       `json:"page_results"`
	EstateAge      int       `json:"estate_age"`
	Square         *rangeInt `json:"square"`
	Price          *rangeInt `json:"price"`
}

type rangeInt struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// MustNewConfig reads config and prepares config struct
func MustNewConfig(path *string) *Config {
	config := new(Config)
	content, err := ioutil.ReadFile(*path)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

    // override options from config file with env vars
	if os.Getenv("DB_HOST") != "" {
		config.Db.Host = os.Getenv("DB_HOST")
	}

	if os.Getenv("DB_USER") != "" {
		config.Db.Username = os.Getenv("DB_USER")
	}

	if os.Getenv("DB_PASSWORD") != "" {
		config.Db.Password = os.Getenv("DB_PASSWORD")
	}

	if os.Getenv("DB_NAME") != "" {
		config.Db.Database = os.Getenv("DB_NAME")
	}

	return config
}

// GetDb returns config for DB
func (c *Config) GetDb() *Db {
	return c.Db
}

// GetEmail returns config for SMTP
func (c *Config) GetEmail() *Email {
	return c.Email
}

// GetSreality returns config for service request to get the list of adverts
func (c *Config) GetSreality() *Sreality {
	return c.Sreality
}
