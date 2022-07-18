package demo

import (
	"crypto/ecdsa"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func StrToPK(s string) (privateKey *ecdsa.PrivateKey, address common.Address) {
	privateKey, err := crypto.HexToECDSA(s)

	if err != nil {
		log.Panicln(err.Error())
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}
