package config

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
)

type Config struct {
    Db *Db
    Email *Email
    Sreality *Sreality
}

type Db struct {
    Engine string `json:"engine"`
    Host string `json:"host"`
    Port int `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type Email struct{
    To string `json:"to"`
    From string `json:"from"`
    Server string `json:"server"`
    Tls bool `json:"tls"`
    TlsPort int `json:"tls_port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Sreality struct{
    Url string `json:"url"`
    UrlDetail string `json:"url_detail"`
    RealityType int `json:"reality_type"`
    OperationType int `json:"operation_type"`
    RealityOptions []int `json:"reality_options"`
    Country int `json:"country"`
    Region []int `json:"region"`
    District []int `json:"district"`
    PageResults int `json:"page_results"`
    EstateAge int `json:"estate_age"`
    Square *rangeInt `json:"square"`
    Price *rangeInt `json:"price"`
}

type rangeInt struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

func MustNewConfig(path *string) *Config {
    config := new(Config);
    content, err := ioutil.ReadFile(*path);
    if err != nil {
        panic(fmt.Sprintf("%v", err))
    }

    err = json.Unmarshal(content, config)
    if err != nil {
        panic(fmt.Sprintf("%v", err))
    }

    return config
}

func (c *Config) GetDb() *Db {
    return c.Db
}

func (c *Config) GetEmail() *Email {
    return c.Email
}

func (c *Config) GetSreality() *Sreality {
    return c.Sreality
}
