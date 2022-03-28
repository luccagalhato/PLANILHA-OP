package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

//Yml ...
var Yml struct {
	API struct {
		Host string `yaml:"host,omitempty"`
		Port string `yaml:"port,omitempty"`
	} `yaml:"api,omitempty"`
	SQL struct {
		Host     string `yaml:"host,omitempty"`
		Port     string `yaml:"port,omitempty"`
		User     string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
	} `yaml:"sql,omitempty"`
	SQLLinx struct {
		Host     string `yaml:"host,omitempty"`
		Port     string `yaml:"port,omitempty"`
		User     string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
	} `yaml:"sqllinx,omitempty"`

	AUTH struct {
		Server string `yaml:"server,omitempty"`
		Port   int    `yaml:"port,omitempty"`
		BaseDN string `yaml:"basedn,omitempty"`
		Grupo  string `yaml:"grupo,omitempty`
	} `yaml:"auth,omitempty"`
}

//LoadConfig ...
func LoadConfig() error {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, &Yml)
}

//CreateConfigFile ...
func CreateConfigFile() {
	if _, err := os.Stat("config.yaml"); err == nil {
		fmt.Println("the 'config.yaml' already exists, do you really want to overwrite? (y/N)")
		var rsp string
		fmt.Scan(&rsp)
		if strings.ToLower(rsp) == "y" {
			writeFile()
		}
		return
	}
	writeFile()
}

func writeFile() {
	b, _ := yaml.Marshal(Yml)
	ioutil.WriteFile("config.yaml", b, 0766)
}
