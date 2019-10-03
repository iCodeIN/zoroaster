package aws

import (
	"log"
	"testing"
	"zoroaster/config"
	"zoroaster/trigger"
)

var psqlClient = PostgresClient{}
var zconf = config.Load("../config")

func init() {
	if zconf.Stage != "DEV" {
		log.Fatal("$STAGE must be DEV to run db tests")
	}
	psqlClient.InitDB(zconf)
}

func TestPostgresClient_All(t *testing.T) {
	// TODO figure out how Go does teardown so I can split these tests;
	// for now I can't be bothered and I'll fit everything in one test,
	// closing the connection only once, at the end.

	// Also note that these tests aren't asserting anything at the moment.
	// The way I'm using them is to run them as a stand-alone module and see
	// if they log any error (they shouldn't).
	// In the future it would be nice to have some real assertions,
	// perhaps with a system that populates from scratch a table.

	defer psqlClient.Close()

	// Log Contract Match
	m := trigger.CnMatch{1, 8888, 10, 0, "matched values", "all values", 1554828248}
	psqlClient.LogCnMatch(m)

	// Update Matching Triggers
	psqlClient.UpdateMatchingTriggers([]int{21, 31})

	// Update Non-Matching Triggers
	psqlClient.UpdateNonMatchingTriggers([]int{21, 31})

	// Log Outcomes
	o1 := trigger.Outcome{"TX outcome", "TX payload"}
	o2 := trigger.Outcome{"CN outcome", "CN payload"}
	psqlClient.LogOutcome(&o1, 1, "wat")
	psqlClient.LogOutcome(&o2, 1, "wac")

	// Load all the active triggers
	_, err := psqlClient.LoadTriggersFromDB("WatchTransactions")
	if err != nil {
		t.Error(err)
	}

}
