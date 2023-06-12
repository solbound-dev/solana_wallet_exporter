package main

import (
	"fmt"
	"math"
	"strconv"
)

const RPC_URL = ""

func GetAccountSolanaBalance(rpc RPC, walletAddress string) (float64, error) {
	balance, err := rpc.GetBalance(walletAddress)
	if err != nil {
		return 0, err
	}

	return float64(balance.Result.Value) / math.Pow(10, 9), nil
}

type Token struct {
	Address string
	Balance float64
}

func GetAccountTokens(rpc RPC, walletAddress string) ([]Token, error) {
	tokenAccounts, err := rpc.GetTokenAccountsByOwner(walletAddress)
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

		amount, err := strconv.ParseFloat(tokenAccountInfo.TokenAmount.Amount, 32)
		if err != nil {
			return nil, err
		}

		balance := amount / math.Pow(10, float64(tokenAccountInfo.TokenAmount.Decimals))

		tokens = append(tokens, Token{Address: tokenAccountInfo.Mint, Balance: balance})
	}

	return tokens, nil
}

func main() {
	rpc := RPC{URL: RPC_URL}

	tokens, err := GetAccountTokens(rpc, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(tokens)

	balance, err := GetAccountSolanaBalance(rpc, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)

}
