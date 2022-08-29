package config

import (
	"encoding/json"
	"github.com/st2projects/ssh-sentinel-client/helpers"
	"os"
)

type ConfigType struct {
	EndPoint   string   `json:"endPoint"`
	APIKey     string   `json:"apiKey"`
	Username   string   `json:"username"`
	Principals []string `json:"principals"`
	PublicKey  string   `json:"publicKey"`
	CertFile   string   `json:"certFile"`
}

var Config *ConfigType

func MakeConfig(configFile string) {
	if !helpers.PathExists(configFile) {
		panic("config file " + configFile + " does not exits")
	}

	configString, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configString, &Config)
	if err != nil {
		panic(err)
	}

}

func (c *ConfigType) GetPublicKey() string {
	return helpers.ExpandPath(c.PublicKey)
}

func (c *ConfigType) GetCertFile() string {
	return helpers.ExpandPath(c.CertFile)
}
