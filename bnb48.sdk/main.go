package bnb48

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

// NewClient creates a client that uses the given RPC client.
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

// ChainID retrieves the current chain ID for transaction replay protection.
func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := ec.puissant.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (ec *Client) SendPuissant(ctx context.Context, txs []string, maxTimestamp int64, acceptReverting []string) (res interface{}, err error) {
	type SendPuissantArgs struct {
		Txs             []string `json:"txs"`
		MaxTimestamp    int64    `json:"maxTimestamp"`
		AcceptReverting []string `json:"acceptReverting"`
	}
	args := SendPuissantArgs{
		Txs:             txs,
		MaxTimestamp:    maxTimestamp,
		AcceptReverting: acceptReverting,
	}
	err = ec.puissant.CallContext(ctx, &res, "eth_sendPuissant", args)
	return
}

// SendPrivateRawTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *Client) SendPrivateRawTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return ec.puissant.CallContext(ctx, nil, "eth_sendPrivateRawTransaction", hexutil.Encode(data))
}

func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	return ec.c.CallContext(ctx, nil, "eth_sendRawTransaction", hexutil.Encode(data))
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}
