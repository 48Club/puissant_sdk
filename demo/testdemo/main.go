package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/bnb48club/puissant_sdk/demo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	var signedTxs []*types.Transaction
	for i := uint64(0); i < 10; i++ {
		tx := types.NewTransaction(nonce+i, fromAddress, value, gasLimit, gasPrice, nil)
		signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		signedTxs = append(signedTxs, signedTx)
	}
	log.Println("signedTxs builded:")

	for i := 0; i < 10; i++ {
		log.Print(signedTxs[i].Hash().String())
	}
	S := time.Now().UnixNano() / 1e3
	var ss []int64
	var ee []int64
	var ress []common.Hash
	var errs []error
	for i := 0; i < 10; i++ {
		ss = append(ss, time.Now().UnixNano()/1e3)
		res, err := client.SendPrivateRawTransaction(context.Background(), signedTxs[0])
		ress = append(ress, res)
		errs = append(errs, err)
		ee = append(ee, time.Now().UnixNano()/1e3)
	}
	for i := 0; i < 10; i++ {
		e, s, res, err := ee[i], ss[i], ress[i], errs[i]
		log.Printf("send tx %d, end time %d, passed %d, res: %s, err: %v", i+1, e, e-s, res.String(), err)

		log.Printf("send tx %d, begin time %d", i+1, s)
	}

	E := time.Now().UnixNano() / 1e3
	log.Printf("send txs, begin time %d, end time %d, passed %d", S, E, E-S)
}
