package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RPC struct {
	URL string
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

func (rpc *RPC) GetBalance(walletAddress string) (rpcGetBalance, error) {
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
				Executable bool   `json:"executable"`
				Lamports   int    `json:"lamports"`
				Owner      string `json:"owner"`
				RentEpoch  int    `json:"rentEpoch"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}

func (rpc *RPC) GetTokenAccountsByOwner(walletAddress string) (rpcGetTokenAccountsByOwnerResp, error) {
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
