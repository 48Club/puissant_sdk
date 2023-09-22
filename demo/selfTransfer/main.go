package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/bnb48club/puissant_sdk/demo"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {

	conf, client := demo.GetClient()

	chainID, err := client.General.ChainID(context.Background())
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Printf("chainID: %s", chainID.String())

	privateKey, fromAddress := demo.StrToPK(conf.Wallet[0])

	nonce, err := client.General.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Panicln(err.Error())
	}

	value := big.NewInt(1e17)
	gasLimit := uint64(21000)

	gasPrice, _ := client.SuggestGasPrice(context.Background())

	var rawTxs []hexutil.Bytes
	// var txs []*types.Transaction
	for k := range make([]int, 2) {
		tx := types.NewTransaction(nonce+uint64(k), fromAddress, value, gasLimit, gasPrice, nil)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Panicln(err.Error())
		}
		// txs = append(txs, signedTx)
		rawTxBytes, _ := rlp.EncodeToBytes(signedTx)

		rawTxs = append(rawTxs, rawTxBytes)
	}

	// send puissant tx
	res, err := client.SendPuissant(context.Background(), rawTxs, uint64(time.Now().Unix()+60), nil)
	// res, err := client.SendPuissantTxs(context.Background(), txs, time.Now().Unix()+60, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println(res)
}
