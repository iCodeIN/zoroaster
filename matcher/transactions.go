package matcher

import (
	"github.com/HAL-xyz/ethrpc"
	"github.com/HAL-xyz/zoroaster/aws"
	"github.com/HAL-xyz/zoroaster/tokenapi"
	"github.com/HAL-xyz/zoroaster/trigger"
	log "github.com/sirupsen/logrus"
	"time"
)

func TxMatcher(blocksChan chan *ethrpc.Block, matchesChan chan trigger.IMatch, idb aws.IDB, api tokenapi.ITokenAPI) {

	for {
		block := <-blocksChan
		start := time.Now()
		log.Info("TX: new -> ", block.Number)

		triggers, err := idb.LoadTriggersFromDB(trigger.WaT)
		if err != nil {
			log.Fatal(err)
		}
		for _, tg := range triggers {
			matchingTxs := trigger.MatchTransaction(tg, block, api)
			for _, m := range matchingTxs {
				matchUUID, err := idb.LogMatch(m)
				if err != nil {
					log.Fatal(err)
				}
				m.MatchUUID = matchUUID
				matchesChan <- *m
			}
		}
		err = idb.SetLastBlockProcessed(block.Number, trigger.WaT)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("\tTX: Processed %d triggers in %s from block %d", len(triggers), time.Since(start), block.Number)
		err = idb.LogAnalytics(trigger.WaT, block.Number, len(triggers), block.Timestamp, start, time.Now())
		if err != nil {
			log.Warn("cannot log analytics: ", err)
		}
	}
}
