package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/hrharder/0x-api-playground/client"
)

func main() {
	baseURL := flag.String("raw-url", "https://kovan.api.0x.org", "set the base 0x API URL")
	ethereumURL := flag.String("ethereum-url", "https://kovan.infura.io/", "set the Ethereum JSONRPC URL (testnet only)")
	sellToken := flag.String("sell-token", "WETH", "set the ticker or address for the sell token")
	buyToken := flag.String("buy-token", "DAI", "set the ticker or address for the buy token")
	sellTokenSize := flag.String("sell-token-size", "", "set the size of the trade in units of the sell token")
	buyTokenSize := flag.String("buy-token-size", "", "set the size of the trade in units of the buy token")
	privateKey := flag.String("private-key", "", "un-prefixed (no '0x') hex-encoded account private key")
	flag.Parse()

	if *privateKey == "" {
		log.Fatal("must set private key (testnet) to use test order filler")
	}

	// size will be nil if empty string is taken in
	buySize, _ := new(big.Int).SetString(*buyTokenSize, 10)
	sellSize, _ := new(big.Int).SetString(*sellTokenSize, 10)
	if buySize == nil && sellSize == nil {
		log.Fatal("one of either buy size or sell size must be specified")
	}
	if buySize != nil && sellSize != nil {
		log.Fatal("must not specify both sell size and buy size; set only one")
	}

	// setup ethereum and 0x-api clients
	zrx, err := client.New(*baseURL)
	if err != nil {
		log.Fatal(err)
	}
	eth, err := ethclient.Dial(*ethereumURL)
	if err != nil {
		log.Fatal(err)
	}

	// call and parse the swap/v0/quote endpoint response
	quote, err := zrx.Quote(*sellToken, *buyToken, sellSize, buySize)
	if err != nil {
		log.Fatal(err)
	}

	// prep key, nonce, and from address
	key, err := crypto.HexToECDSA(*privateKey)
	opts := bind.NewKeyedTransactor(key)
	nonce, err := eth.NonceAt(context.Background(), opts.From, nil)

	// estimate gas limit for transaction
	gasLimit, err := eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     opts.From,
		To:       &quote.To,
		GasPrice: quote.GasPrice,
		Data:     quote.Data,
		Value:    quote.Value,
	})
	if err != nil {
		log.Fatal(err)
	}

	// construct un-signed tx with values from above and quote
	tx := types.NewTransaction(nonce, quote.To, quote.Value, gasLimit, quote.GasPrice, quote.Data)

	// sign tx and send
	chainID, err := eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signer := types.NewEIP155Signer(chainID)
	signedTx, _ := types.SignTx(tx, signer, key)
	if err := eth.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("fill sent; transaction hash: %s", signedTx.Hash().Hex())
}
