package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/du5/puissant_demo/bnb48.sdk"
	"github.com/du5/puissant_demo/demo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	conf := demo.GetConf("config.yaml")

	client, err := bnb48.Dial("https://testnet-fonce-bsc.bnb48.club", "https://testnet-puissant-bsc.bnb48.club")
	if err != nil {
		log.Panicln(err.Error())
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Printf("chainID: %s", chainID.String())

	countWallet := len(conf.Wallet)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	value := big.NewInt(2e18)

	var rawTxs []string
	// var txs []*types.Transaction
	for i := 0; i < countWallet; i++ {
		pk := conf.Wallet[i]
		privateKey, fromAddress := demo.StrToPK(pk)
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Panicln(err.Error())
		}

		gasLimit := uint64(21000)
		toPk := conf.Wallet[0]
		if i < countWallet-1 {
			toPk = conf.Wallet[i+1]
		}
		_, toAddress := demo.StrToPK(toPk)
		// gas sort
		thisGas := big.NewInt(0).Add(big.NewInt(int64(countWallet-i)), gasPrice)
		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, thisGas, nil)
		// change next value
		value.Sub(value, thisGas.Mul(thisGas, big.NewInt(int64(gasLimit))))
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Panicln(err.Error())
		}
		// txs = append(txs, signedTx)
		rawTxBytes, _ := rlp.EncodeToBytes(types.Transactions{signedTx}[0])
		rawTxHex := hexutil.Encode(rawTxBytes)

		rawTxs = append(rawTxs, rawTxHex)
	}

	// send puissant tx
	res, err := client.SendPuissant(context.Background(), rawTxs, time.Now().Unix()+60, nil)
	// res, err := client.SendPuissantTxs(context.Background(), txs, time.Now().Unix()+60, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println(res)
}
