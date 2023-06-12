package main

import (
	"fmt"
)

const RPC_URL = ""

func main() {
	tokens, err := GetAccountTokens(RPC_URL, "")
	fmt.Println(tokens)
	fmt.Println(err)
	fmt.Println("hi")

}
