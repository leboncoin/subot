package analytics

import (
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleDeletedMessage godoc
// @Summary Handles the deletion of a message
// @Description If the message that was deleted only had one answer
// @Description which came from the bot, this endpoint returns an action
// @Description to delete this reply to the message in order to delete
// @Description the original message from the main thread.
// @Description It will change the status of the message to deleted in the storage.
// @Tags Analytics
// @ID handle-deleted-message
// @Accept  json
// @Produce  json
// @Param message path int true "Account ID"
// @Router /analytics/deleted [post]
func (a Analyser) HandleDeletedMessage(message globals.Message) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	reply.Action = globals.Nothing
	log.WithFields(log.Fields{"message": message}).Debug("Handle deleted message")

	log.Debug("Check if message has responses")
	originalTs := message.DeletedTs
	if originalTs == "" {
		log.Debug("original ts is not deleted ts")
		originalTs = message.EditedTs
	}
	storedMessages, err := a.ESClient.QueryRangeMessages(originalTs, originalTs)
	log.WithFields(log.Fields{"messages": storedMessages}).Debug("Found stored message for this timestamp")
	if err != nil {
		log.Error("Failed to fetch last messages", err)
	}
	if len(storedMessages) == 0 {
		log.Error("Found no messages")
		return
	}
	log.Debug("Get replies for message")
	deletedMessageID := storedMessages[0].ID
	deletedMessage := storedMessages[0]
	log.WithFields(log.Fields{"message": deletedMessage}).Debug("Debug message")
	if len(deletedMessage.Replies) == 1 {
		log.Debug("message has only one response")
		if deletedMessage.Replies[0].FromBot {
			log.Debug("Response is from me")
			reply.Action = globals.DeleteMessage
			reply.Ts = deletedMessage.Replies[0].Timestamp
		}
	}
	deletedMessage.Type = "deleted"
	deletedMessage.RemindAt = ""
	err = a.ESClient.AddMessage(deletedMessage, deletedMessageID)
	if err != nil {
		log.Error("Error while changing message type to deleted :", err)
		return
	}

	return []globals.SlackResponse{reply}, nil
}
