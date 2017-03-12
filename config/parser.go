package config

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
)

//const path = "./config.json"

type Config struct {
    Db *Db
    Email *Email
    Sreality *Sreality
}

type Db struct {
    Engine string
    Host string
    Port int
    Username string
    Password string
    Database string
}

type Email struct{
    To string
    From string
    Server string
    Tls bool
    TlsPort int `json:"tls_port"`
    Username string
    Password string
}

type Sreality struct{
    Url string
    UrlDetail string `json:"url_detail"`
    RealityType int `json:"reality_type"`
    OperationType int `json:"operation_type"`
    RealityOptions []int `json:"reality_options"`
    Country int
    Region []int
    District []int
    PageResults int `json:"page_results"`
    EstateAge int `json:"estate_age"`
    Square *rangeInt
    Price *rangeInt
}

type rangeInt struct {
    Min int
    Max int
}

func NewConfig(path *string) (*Config) {
    config := new(Config);
    content, err := ioutil.ReadFile(*path);
    err = json.Unmarshal(content, config)
    if err != nil {
        panic(fmt.Sprintf("%v", err))
    }

    return config
}

func (c *Config) GetDb() (*Db) {
    return c.Db
}

func (c *Config) GetEmail() (*Email) {
    return c.Email
}

func (c *Config) GetSreality() (*Sreality) {
    return c.Sreality
}
