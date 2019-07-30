package main

import (
	"github.com/onrik/ethrpc"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
	"zoroaster/aws"
	"zoroaster/config"
	"zoroaster/eth"
	"zoroaster/matcher"
)

func main() {

	// Load Config
	zconf := config.Load("config")

	// Load AWS SES session
	sesSession := aws.GetSESSession()

	// Persist logs
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})
	log.SetLevel(log.DebugLevel)
	f, err := os.OpenFile(zconf.LogsFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Init Postgres DB client
	psqlClient := aws.PostgresClient{}
	psqlClient.InitDB(zconf)

	// HTTP client
	httpClient := http.Client{}

	// ETH client
	ethClient := ethrpc.New(zconf.EthNode)

	// Channels are buffered because a) contracts is slower, and b) so I can run wac/wat independently for tests
	txBlocksChan := make(chan *ethrpc.Block, 100)
	contractsBlocksChan := make(chan *ethrpc.Block, 100)
	matchesChan := make(chan interface{})

	// Poll ETH node
	go eth.BlocksPoller(txBlocksChan, contractsBlocksChan, ethClient, &psqlClient)

	// Watch a Transaction
	go matcher.TxMatcher(txBlocksChan, matchesChan, &psqlClient)

	// Watch a Contract
	go matcher.ContractMatcher(contractsBlocksChan, matchesChan, eth.GetModifiedAccounts, &psqlClient, ethClient)

	// Main routine - process matches
	for {
		match := <-matchesChan
		go matcher.ProcessMatch(match, &psqlClient, sesSession, &httpClient)
	}
}
