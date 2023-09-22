package demo

import (
	"crypto/ecdsa"
	"log"
	"os"

	bnb48_sdk "github.com/bnb48club/puissant_sdk/bnb48.sdk"
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
func getConf(p string) (conf *Conf) {
	yamlFile, err := os.ReadFile(p)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Panicln(err.Error())
	}
	return
}

func GetClient() (conf *Conf, client *bnb48_sdk.Client) {
	conf = getConf("config.yaml")

	client, err := bnb48_sdk.Dial("https://1gwei.48.club", os.Getenv("RPC"))
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
