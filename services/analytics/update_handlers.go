package analytics

import (
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleUpdatedMessage godoc
// @Summary Handles messages that are updated
// @Description Returns an empty replies and stores the new message content to storage.
// @Description It doesn't analyse the message again
// @Tags Analytics
// @ID handle-updated-message
// @Accept  json
// @Produce  json
// @Param message body object true "The content of the updated message"
// @Router /analytics/updated [post]
func (a Analyser) HandleUpdatedMessage(message globals.Message) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	reply.Action = globals.Nothing
	log.WithFields(log.Fields{"message": message}).Debug("Handle updated message")

	log.Debug("Look for message to update")
	originalTs := message.Timestamp
	storedMessages, err := a.ESClient.QueryRangeMessages(originalTs, originalTs)
	log.WithFields(log.Fields{"messages": storedMessages}).Debug("Found stored message for this timestamp")
	if err != nil {
		log.Error("Failed to fetch last messages", err)
		return []globals.SlackResponse{reply}, err
	}
	if len(storedMessages) == 0 {
		log.Error("Found no messages")
		return []globals.SlackResponse{reply}, nil
	}
	log.Debug("Get replies for message")
	updatedMessageID := storedMessages[0].ID
	updatedMessage := storedMessages[0]
	log.WithFields(log.Fields{"message": updatedMessage}).Debug("Debug message")
	if updatedMessage.Text != message.Text {
		updatedMessage.Text = message.Text
		err = a.ESClient.AddMessage(updatedMessage, updatedMessageID)
		if err != nil {
			log.Error("Error while changing message text :", err)
			return []globals.SlackResponse{reply}, err
		}
	}

	return []globals.SlackResponse{reply}, nil
}
