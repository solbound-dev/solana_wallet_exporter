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
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	solanaBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solana_wallet_balance",
			Help: "Balance of a given wallet in SOL.",
		},
		[]string{"wallet"},
	)
	tokenBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "solana_wallet_token_balance",
			Help: "Balance of a given wallet in a given token.",
		},
		[]string{"wallet", "token"},
	)
	lastUpdate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "solana_wallet_last_update_ts",
			Help: "Timestamp of the last update.",
		},
	)
)

func recordMetrics(rpc RPC, wallets string, cacheSeconds int) {
	go func() {
		for {
			for _, wallet := range strings.Split(wallets, ",") {
				balance, err := rpc.GetAccountSolanaBalance(wallet)
				if err != nil {
					log.Println(err)
				}
				solanaBalance.WithLabelValues(wallet).Set(balance)

				tokens, err := rpc.GetAccountTokens(wallet)
				if err != nil {
					log.Println(err)
				}
				for _, token := range tokens {
					tokenBalance.WithLabelValues(wallet, token.Address).Set(token.Balance)
				}
			}

			lastUpdate.Set(float64(time.Now().Unix()))
			time.Sleep(time.Duration(cacheSeconds) * time.Second)
		}
	}()
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":18899").String()
		rpcURL        = kingpin.Flag("solana.rpc", "Solana RPC provider URL.").Required().String()
		wallets       = kingpin.Flag("solana.wallets", "Comma separated list of solana wallets.").Required().String()
		cacheSeconds  = kingpin.Flag("solana.cacheseconds", "Number of seconds to cache values for.").Default("300").Int()
	)
	kingpin.Parse()

	recordMetrics(RPC{URL: *rpcURL}, *wallets, *cacheSeconds)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on address %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server (%s)", err)
	}
}
