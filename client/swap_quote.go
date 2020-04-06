package client

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (c *Client) Quote(sellToken string, buyToken string, sellAmount *big.Int, buyAmount *big.Int) (*SwapResponse, error) {
	response := new(SwapResponse)
	args := Args{
		"buyToken":  buyToken,
		"sellToken": sellToken,
	}
	if sellAmount != nil {
		args["sellAmount"] = sellAmount.String()
	}
	if buyAmount != nil {
		args["buyAmount"] = buyAmount.String()
	}

	rawRes, err := c.get("swap", 0, "quote", args)
	if err != nil {
		return nil, fmt.Errorf("(client) failed to fetch quote: %w", err)
	}

	if err := json.NewDecoder(rawRes.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("(client) failed to decode quote: %w", err)
	}

	return response, nil
}

type Source struct {
	Name       string
	Proportion *big.Float
}

type sourceJSON struct {
	Name       string `json:"name"`
	Proportion string `json:"proportion"`
}

func (src *Source) UnmarshalJSON(data []byte) error {
	rawSource := &sourceJSON{}
	if err := json.Unmarshal(data, rawSource); err != nil {
		return fmt.Errorf("(client) failed to unmarshal: %w", err)
	}

	parsedProportion, ok := new(big.Float).SetString(rawSource.Proportion)
	if !ok {
		return fmt.Errorf("(client) failed to parse float proportion")
	}

	src.Name = rawSource.Name
	src.Proportion = parsedProportion
	return nil
}

type SwapResponse struct {
	Price            *big.Float
	To               common.Address
	Data             []byte
	Value            *big.Int
	GasPrice         *big.Int
	Gas              *big.Int
	ProtocolFee      *big.Int
	BuyAmount        *big.Int
	SellAmount       *big.Int
	Sources          []*Source
	BuyTokenAddress  common.Address
	SellTokenAddress common.Address
}

type swapResponseJSON struct {
	Price            string        `json:"price"`
	To               string        `json:"to"`
	Data             hexutil.Bytes `json:"data"`
	Value            string        `json:"value"`
	GasPrice         string        `json:"gasPrice"`
	Gas              string        `json:"gas"`
	ProtocolFee      string        `json:"protocolFee"`
	BuyAmount        string        `json:"buyAmount"`
	SellAmount       string        `json:"sellAmount"`
	Sources          []*Source     `json:"sources"`
	BuyTokenAddress  string        `json:"buyTokenAddress"`
	SellTokenAddress string        `json:"sellTokenAddress"`
}

func (sr *SwapResponse) UnmarshalJSON(data []byte) error {
	rawResponse := &swapResponseJSON{}
	if err := json.Unmarshal(data, rawResponse); err != nil {
		return fmt.Errorf("(client) failed to unmarshal: %w", err)
	}

	price, ok := new(big.Float).SetString(rawResponse.Price)
	if !ok {
		return fmt.Errorf("(client) failed to parse float price")
	}
	value, ok := new(big.Int).SetString(rawResponse.Value, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer value")
	}
	gasPrice, ok := new(big.Int).SetString(rawResponse.GasPrice, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer gasPrice")
	}
	gas, ok := new(big.Int).SetString(rawResponse.Gas, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer gas")
	}
	protocolFee, ok := new(big.Int).SetString(rawResponse.ProtocolFee, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer protocolFee")
	}
	buyAmount, ok := new(big.Int).SetString(rawResponse.BuyAmount, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer buyAmount")
	}
	sellAmount, ok := new(big.Int).SetString(rawResponse.SellAmount, 10)
	if !ok {
		return fmt.Errorf("(client) failed to parse integer sellAmount")
	}

	sr.Price = price
	sr.To = common.HexToAddress(rawResponse.To)
	sr.Data = []byte(rawResponse.Data)
	sr.Value = value
	sr.GasPrice = gasPrice
	sr.Gas = gas
	sr.ProtocolFee = protocolFee
	sr.BuyAmount = buyAmount
	sr.SellAmount = sellAmount
	sr.Sources = rawResponse.Sources
	sr.BuyTokenAddress = common.HexToAddress(rawResponse.BuyTokenAddress)
	sr.SellTokenAddress = common.HexToAddress(rawResponse.SellTokenAddress)
	return nil
}
