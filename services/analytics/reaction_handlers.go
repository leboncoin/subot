package analytics

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)


// HandleReaction godoc
// @Summary Handles reactions, looking for a heavy check mark.
// @Description Returns an empty reply but stores the new status
// @Description if the reaction is heavy_check_mark.
// @Description It also calculates response time
// @Description based on local time and message timestamp.
// @Tags Analytics
// @ID handle-reaction
// @Accept  json
// @Produce  json
// @Param reaction body object true "Content of the reaction"
// @Router /analytics/reaction [post]
func (a Analyser) HandleReaction(reaction globals.Reaction) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	log.Debug("Get original message")
	originalMessages, err := a.ESClient.QueryRangeMessages(reaction.MessageTs, reaction.MessageTs)
	if len(originalMessages) == 0 {
		log.Debug("Original message not found, return")
		return
	}

	originalMessage := originalMessages[0]
	originalMessage.Reactions = append(originalMessage.Reactions, reaction)
	log.Debug("Calculate and save resolution time if reaction is :heavy_check_mark:")

	if reaction.Name == "heavy_check_mark" && originalMessage.Status != "fixed" {
		resolutionTime := (globals.ParseDuration(reaction.Timestamp) - globals.ParseDuration(reaction.MessageTs)) / 60
		originalMessage.ResolutionTime = time.Duration(resolutionTime)
		originalMessage.Status = "fixed"
		originalMessage.RemindAt = ""
	}
	log.WithFields(log.Fields{"event": reaction}).Debug("Save reaction for message")
	err = a.ESClient.AddMessage(originalMessage, originalMessages[0].ID)

	return []globals.SlackResponse{reply}, nil
}
