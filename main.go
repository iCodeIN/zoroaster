package main

import (
	"github.com/onrik/ethrpc"
	"log"
	"os"
	"time"
	"zoroaster/aws"
	"zoroaster/rpc"
	"zoroaster/triggers"
)

func main() {

	// Load table config
	table := os.Getenv("DB_TABLE")
	if table == "" {
		table = "trigger1"
	}

	// Connect to triggers' DB
	aws.InitDB()

	// Poll ETH node
	c := make(chan *ethrpc.Block)
	go rpc.PollForLastBlock(c)

	// Main routine
	for {
		block := <-c
		start := time.Now()
		log.Println("New block: #", block.Number)

		triggers, err := aws.LoadTriggersFromDB(table)
		if err != nil {
			log.Fatal(err)
		}

		for _, tg := range triggers {
			txs := trigger.MatchTrigger(tg, block)
			for _, tx := range txs {
				log.Printf("\tTrigger %d matched transaction "+
					"https://etherscan.io/tx/%s", tg.TriggerId, tx.Hash)
			}
		}
		log.Printf("\tProcessed %d triggers in %s", len(triggers), time.Since(start))
	}

}