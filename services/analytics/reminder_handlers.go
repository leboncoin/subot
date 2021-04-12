package analytics

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleRemindersRequest godoc
// @Summary Checks for reminders
// @Description Returns the list of reminders to send
// @Description when last reply on a message is an hour old
// @Tags Analytics
// @ID handle-reminders-request
// @Produce  json
// @Router /analytics/reminders [post]
func (a Analyser) HandleRemindersRequest() (replies []globals.SlackResponse, err error) {
	messages, err := a.ESClient.QueryReminderMessages()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Got an error while fetching messages with reminders")
		return
	}

	for _, message := range messages {
		var reply globals.SlackResponse
		reply.Action = globals.ReplyMessage
		reply.Text = fmt.Sprintf("Du nouveau <@%s> ?", a.getFiremanID())
		reply.Ts = message.Timestamp
		replies = append(replies, reply)
	}

	return
}
