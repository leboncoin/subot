package analytics

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleReplies godoc
// @Summary Handles message replies and calculates response time
// @Description returns an empty reply but stores the replies to the storage.
// @Description Calculates response times only for the first team member answer
// @Tags Analytics
// @ID handle-replies
// @Accept  json
// @Produce  json
// @Param message body object true "The reply that was posted"
// @Router /analytics/thread [post]
func (a Analyser) HandleReplies(message globals.Reply) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	log.Debug("Get original message")
	originalMessages, err := a.ESClient.QueryRangeMessages(message.ThreadTs, message.ThreadTs)
	if len(originalMessages) == 0 {
		log.Debug("Original message not found, return")
		return
	}
	originalMessage := originalMessages[0]
	originalMessage.Replies = append(originalMessage.Replies, message)
	if !message.FromBot && originalMessage.Status != "fixed" {
		originalMessage.RemindAt = strconv.FormatInt(time.Now().Add(reminderInterval).Unix(), 10)
	} else {
		originalMessage.RemindAt = ""
	}

	log.Debug("Calculate and save response time if reply owner is team member")

	isTeamMessage, err := a.ESClient.IsTeamMember(message.UserID)
	if err != nil {
		return
	}
	if isTeamMessage && originalMessage.Status == "unresponded" {
		log.Debug("Its a team message")
		log.Debug("Set responded status")
		originalMessage.Status = "responded"
		responseTime := (globals.ParseDuration(message.Timestamp) - globals.ParseDuration(message.ThreadTs)) / 60
		log.Debug("Response time is : ", responseTime)
		originalMessage.ResponseTime = time.Duration(responseTime)
	}

	log.WithFields(log.Fields{"event": message}).Debug("Save reply for message")
	err = a.ESClient.AddMessage(originalMessage, originalMessages[0].ID)

	return []globals.SlackResponse{reply}, nil
}
