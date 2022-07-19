package bnb48

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c        *rpc.Client
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
		c:        cs[0],
		puissant: cs[1],
	}
}

func (ec *Client) Close() {
	ec.c.Close()
	ec.puissant.Close()
}

func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return ethclient.NewClient(ec.c).ChainID(ctx)
}

func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return ethclient.NewClient(ec.puissant).SuggestGasPrice(ctx)
}

func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return ethclient.NewClient(ec.c).SendTransaction(ctx, tx)
}

func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return ethclient.NewClient(ec.c).PendingNonceAt(ctx, account)
}

type SendPuissantArgs struct {
	Txs             []string `json:"txs"`
	MaxTimestamp    int64    `json:"maxTimestamp"`
	AcceptReverting []string `json:"acceptReverting"`
}

func (ec *Client) SendPuissant(ctx context.Context, txs []string, maxTimestamp int64, acceptReverting []string) (res interface{}, err error) {
	err = ec.puissant.CallContext(ctx, &res, "eth_sendPuissant", SendPuissantArgs{
		Txs:             txs,
		MaxTimestamp:    maxTimestamp,
		AcceptReverting: acceptReverting,
	})
	return
}

func (ec *Client) SendPuissantTxs(ctx context.Context, txs []*types.Transaction, maxTimestamp int64, acceptReverting []*types.Transaction) (interface{}, error) {
	signedRawTxs := make([][]string, 2)
	for k, signedTxs := range [][]*types.Transaction{txs, acceptReverting} {
		for _, signedTx := range signedTxs {
			rawTxBytes, err := rlp.EncodeToBytes(types.Transactions{signedTx}[0])
			if err != nil {
				return nil, err
			}
			signedRawTxs[k] = append(signedRawTxs[k], hexutil.Encode(rawTxBytes))
		}
	}

	return ec.SendPuissant(ctx, signedRawTxs[0], maxTimestamp, signedRawTxs[1])
}

func (ec *Client) SendPrivateRawTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return ec.puissant.CallContext(ctx, nil, "eth_sendPrivateRawTransaction", hexutil.Encode(data))
}
