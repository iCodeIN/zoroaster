package matcher

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	log "github.com/sirupsen/logrus"
	"zoroaster/action"
	"zoroaster/aws"
	"zoroaster/trigger"
)

func ProcessMatch(match trigger.IMatch, idb aws.IDB, iEmail sesiface.SESAPI, httpCli aws.IHttpClient) []*trigger.Outcome {

	acts, err := idb.GetActions(match.GetTriggerUUID(), match.GetUserUUID())
	if err != nil {
		log.Fatalf("cannot get actions from db: %v", err)
	}
	log.Debugf("\tMatched %d actions", len(acts))

	outcomes := action.ProcessActions(acts, match, iEmail, httpCli)
	for _, out := range outcomes {
		idb.LogOutcome(out, match.GetMatchUUID())
		log.Debug("\tLogged outcome for match id ", match.GetMatchUUID())
	}
	return outcomes
}
