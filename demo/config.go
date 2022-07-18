package demo

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Wallet []string `yaml:"wallet"`
}

// GetConf returns a Conf struct from a yaml file
//
// p: path to yaml file
//
func GetConf(p string) (conf *Conf) {
	yamlFile, err := ioutil.ReadFile(p)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Panicln(err.Error())
	}
	return
}
