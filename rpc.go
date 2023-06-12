package main

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Token struct {
	Address string
	Balance float64
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

func GetAccountTokens(rpcURL string, walletAddress string) ([]Token, error) {
	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(`
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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenAccounts rpcGetTokenAccountsByOwnerResp
	if err := json.Unmarshal(body, &tokenAccounts); err != nil {
		return nil, err
	}

	var tokens []Token
	for _, tokenAccount := range tokenAccounts.Result.Value {
		tokenAccountInfo := tokenAccount.Account.Data.Parsed.Info

		// skip NFTs
		if tokenAccountInfo.TokenAmount.Amount == "1" && tokenAccountInfo.TokenAmount.Decimals == 0 {
			continue
		}

		amount, err := strconv.ParseFloat(tokenAccountInfo.TokenAmount.Amount, 32)
		if err != nil {
			return nil, err
		}

		balance := amount / math.Pow(10, float64(tokenAccountInfo.TokenAmount.Decimals))

		tokens = append(tokens, Token{Address: tokenAccountInfo.Mint, Balance: balance})
	}

	return tokens, nil
}
