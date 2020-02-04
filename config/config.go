package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var Conf, err = ReadConfig()

type Config struct {
	App struct {
		Name        string `yaml:"Name"`
		LimitPerMin string `yaml:"LimitPerMin"`
	} `yaml:"App"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"Postgres"`
	JWT struct {
		PrivateBytes string `yaml:"PrivateBytes"`
		PublicBytes  string `yaml:"PublicBytes"`
		ExpireHours  int    `yaml:"ExpireHours"`
	} `yaml:"JWT"`
}

// ReadConfig read config from file
func ReadConfig() (c Config, err error) {
	env := os.Getenv("iris_env")
	if env == "" {
		env = "dev"
	}
	filename := "app." + env + ".yml"
	workpath, err := os.Getwd()
	file, err := ioutil.ReadFile(workpath + "/config/" + filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return
	}
	return
}
