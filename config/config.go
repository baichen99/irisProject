package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var Conf, err = ReadConfig()

type Config struct {
	App struct {
		Name string	`yaml:"Name"`
	} `yaml:"App"`
	JWT struct {
		PrivateBytes string `yaml:"PrivateBytes"`
		PublicBytes  string `yaml:"PublicBytes"`
		ExpireHours	int	`yaml:"ExpireHours"`
	} `yaml:"JWT"`
}


func ReadConfig() (c Config, err error) {
	file, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return
	}
	return
}
