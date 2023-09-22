package bnb48

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	General  *ethclient.Client
	puissant *rpc.Client
}

func Dial(general, puissant string) (*Client, error) {
	return DialContext(context.Background(), general, puissant)
}

func DialContext(ctx context.Context, rpcs ...string) (*Client, error) {
	var cs []*rpc.Client
	for _, rawurl := range rpcs {
		c, err := rpc.DialContext(ctx, rawurl)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}

	return NewClient(cs...), nil
}

func NewClient(cs ...*rpc.Client) *Client {
	return &Client{
		General:  ethclient.NewClient(cs[0]),
		puissant: cs[1],
	}
}

func (ec *Client) Close() {
	ec.General.Close()
	ec.puissant.Close()
}

func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return ethclient.NewClient(ec.puissant).SuggestGasPrice(ctx)
}

type SendPuissantArgs struct {
	Txs             []hexutil.Bytes `json:"txs"`
	MaxTimestamp    uint64          `json:"maxTimestamp"`
	AcceptReverting []common.Hash   `json:"acceptReverting"`
}

func (ec *Client) SendPuissant(ctx context.Context, txs []hexutil.Bytes, maxTimestamp uint64, acceptReverting []common.Hash) (res string, err error) {
	err = ec.puissant.CallContext(ctx, &res, "eth_sendPuissant", SendPuissantArgs{
		Txs:             txs,
		MaxTimestamp:    maxTimestamp,
		AcceptReverting: acceptReverting,
	})
	return
}

func (ec *Client) SendPuissantTxs(ctx context.Context, txs []*types.Transaction, maxTimestamp uint64, acceptReverting []*types.Transaction) (res string, err error) {
	txsBytes := []hexutil.Bytes{}
	for _, signedTx := range txs {
		var rawTxBytes hexutil.Bytes
		rawTxBytes, err = signedTx.MarshalBinary()
		if err != nil {
			return
		}
		txsBytes = append(txsBytes, rawTxBytes)
	}

	txsHash := []common.Hash{}
	for _, v := range acceptReverting {
		txsHash = append(txsHash, v.Hash())
	}

	return ec.SendPuissant(ctx, txsBytes, maxTimestamp, txsHash)
}

func (ec *Client) SendPrivateRawTransaction(ctx context.Context, tx *types.Transaction) (h common.Hash, err error) {
	data, err := tx.MarshalBinary()
	if err != nil {
		return
	}
	err = ec.puissant.CallContext(ctx, &h, "eth_sendPrivateRawTransaction", hexutil.Encode(data))
	return
}
