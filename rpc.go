// Copyright 2023 Bartol Deak
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type RPC struct {
	URL string
}

func (rpc *RPC) GetAccountSolanaBalance(walletAddress string) (float64, error) {
	balance, err := rpc.rpcGetBalance(walletAddress)
	if err != nil {
		return 0, err
	}

	return float64(balance.Result.Value) / math.Pow(10, 9), nil
}

type rpcGetBalance struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			APIVersion string `json:"apiVersion"`
			Slot       int    `json:"slot"`
		} `json:"context"`
		Value int `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

func (rpc *RPC) rpcGetBalance(walletAddress string) (rpcGetBalance, error) {
	resp, err := http.Post(rpc.URL, "application/json", strings.NewReader(`
	{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getBalance",
		"params": [
			"`+walletAddress+`"
		]
	}
	`))
	if err != nil {
		return rpcGetBalance{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rpcGetBalance{}, err
	}

	var data rpcGetBalance
	if err := json.Unmarshal(body, &data); err != nil {
		return rpcGetBalance{}, err
	}

	if data.ID == 0 {
		err := fmt.Errorf("error: %s", string(body))
		return rpcGetBalance{}, err
	}

	return data, nil
}

type Token struct {
	Address string
	Balance float64
}

func (rpc *RPC) GetAccountTokens(walletAddress string) ([]Token, error) {
	tokenAccounts, err := rpc.rpcGetTokenAccountsByOwner(walletAddress)
	if err != nil {
		return nil, err
	}

	var tokens []Token
	for _, tokenAccount := range tokenAccounts.Result.Value {
		tokenAccountInfo := tokenAccount.Account.Data.Parsed.Info

		// skip NFTs
		if tokenAccountInfo.TokenAmount.Amount == "1" && tokenAccountInfo.TokenAmount.Decimals == 0 {
			continue
		}

		amount, err := strconv.ParseFloat(tokenAccountInfo.TokenAmount.Amount, 64)
		if err != nil {
			return nil, err
		}

		balance := amount / math.Pow(10, float64(tokenAccountInfo.TokenAmount.Decimals))

		tokens = append(tokens, Token{Address: tokenAccountInfo.Mint, Balance: balance})
	}

	return tokens, nil
}

type rpcGetTokenAccountsByOwnerResp struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			APIVersion string `json:"apiVersion"`
			Slot       int    `json:"slot"`
		} `json:"context"`
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							IsNative    bool   `json:"isNative"`
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							State       string `json:"state"`
							TokenAmount struct {
								Amount         string  `json:"amount"`
								Decimals       int     `json:"decimals"`
								UIAmount       float64 `json:"uiAmount"`
								UIAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program string `json:"program"`
					Space   int    `json:"space"`
				} `json:"data"`
				Executable bool    `json:"executable"`
				Lamports   int     `json:"lamports"`
				Owner      string  `json:"owner"`
				RentEpoch  big.Int `json:"rentEpoch"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

func (rpc *RPC) rpcGetTokenAccountsByOwner(walletAddress string) (rpcGetTokenAccountsByOwnerResp, error) {
	resp, err := http.Post(rpc.URL, "application/json", strings.NewReader(`
	{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "getTokenAccountsByOwner",
		"params": [
			"`+walletAddress+`",
			{
				"programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
			},
			{
				"encoding": "jsonParsed"
			}
		]
	}
	`))
	if err != nil {
		return rpcGetTokenAccountsByOwnerResp{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rpcGetTokenAccountsByOwnerResp{}, err
	}

	var data rpcGetTokenAccountsByOwnerResp
	if err := json.Unmarshal(body, &data); err != nil {
		return rpcGetTokenAccountsByOwnerResp{}, err
	}

	if data.ID == 0 {
		err := fmt.Errorf("error: %s", string(body))
		return rpcGetTokenAccountsByOwnerResp{}, err
	}

	return data, nil
}
