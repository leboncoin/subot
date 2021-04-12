package replier

import (
	"bytes"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/pkg/slack"
)

// CatchUp godoc
// @Summary Backport messages for the given period
// @Description Curls the slack to retrieve all the messages
// @Description for the period (by batch of 10 messages)
// @Description For every batch, it then calls CatchUpBatch to
// @Description analyse the batch of messages
// @ID catchup
// @Produce  json
// @Param start query string true "Start of the"
// @Param end query string true "End of the period"
// @Router /catchup [get]
func (h Handler) CatchUp(start string, end string) {
	log.Debugf("Catching up messages from %s to %s", start, end)
	cursor := ""
	hasMore := true
	for hasMore {
		s, err := globals.ParseDate(start)
		if err != nil {
			return
		}
		e, err := globals.ParseDate(end)
		if err != nil {
			return
		}
		slackResponse, err := h.Slack.ReadMessages(s, e, cursor, 10)
		if err != nil {
			return
		}
		go h.CatchUpBatch(slackResponse.Messages)

		hasMore = slackResponse.HasMore
		cursor = slackResponse.Metadata.NextCursor
	}
	return
}

// CatchUpBatch runs the data consolidation and message analysis + storage for a batch of events
func (h Handler) CatchUpBatch(events []slack.Event) {
	time.Sleep(time.Duration(1) * time.Minute)
	var messages []globals.Message

	for _, event := range events {
		message := h.Slack.GetMessage(event)
		messages = append(messages, message)
	}

	jsonBody, err := json.Marshal(messages)
	if err != nil {
		log.Error("error while decoding event", err)
		return
	}

	_, err = h.callAnalyticsAPI("POST", "batch", bytes.NewReader(jsonBody))
}
