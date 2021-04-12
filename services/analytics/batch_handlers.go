package analytics

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleBatchMessage godoc
// @Summary Analyses a batch of messages and stores them into the database
// @Description This endpoint returns no reply but saves all the messages
// @Description that are not already present in the storage after being analysed.
// @Tags Analytics
// @ID handle-batch-message
// @Accept  json
// @Produce  json
// @Param message path int true "Account ID"
// @Router /analytics/batch [post]
func (a Analyser) HandleBatchMessage(messages []globals.Message) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	var documentID string
	reply.Action = globals.Nothing
	for _, message := range messages {
		log.Debug("Check if message already exists")
		storedMessages, err := a.ESClient.QueryRangeMessages(message.Timestamp, message.Timestamp)
		if err != nil {
			continue
		}
		if len(storedMessages) > 0 {
			continue
		}
		if message.Type == globals.Join || message.Type == globals.Left {
			continue
		}

		log.Debug("Check if message comes from team member")
		isTeamMessage, err := a.ESClient.IsTeamMember(message.UserID)
		if err != nil {
			continue
		}
		if isTeamMessage {
			if message.Type == "topic" {
				log.Debug("Found topic change, save fireman")
				err = a.ESClient.AddFireman(message)
				if err != nil {
					log.Error("Got an error while saving fireman", err)
				}
				continue
			} else {
				message.Type = "team"
			}
		} else {
			message.Type = "user"
		}

		log.Debug("Get tools for event")
		//incidents, err := a.ESClient.QueryIncidents(tools)
		tools, err := a.ESClient.QueryTools(message.Text)
		if err != nil {
			log.Error("Got an error while querying tools", err)
			continue
		}
		log.Debug("Get labels for event")
		//incidents, err := a.ESClient.QueryIncidents(tools)
		for _, t := range tools {
			log.WithFields(log.Fields{"tool": t}).Debug("Send link to documentation for tool")
			//reply.Text = reply.Text + " La documentation de " + t + " est disponible ici : link to doc"
		}
		// Error while querying labels
		labels, err := a.ESClient.QueryLabels(message.Text)
		if err != nil {
			log.Error("Got an error while querying labels", err)
			continue
		}
		log.WithFields(log.Fields{"tools": tools /*, "labels": labels*/}).Debug("Got tools and labels")
		log.Debug("Save message to elasticsearch")

		message.Tools = tools
		message.Labels = labels
		message.Status = "unresponded"
		if len(message.Replies) > 0 {
			for _, reply := range message.Replies {
				isTeamReply, err := a.ESClient.IsTeamMember(reply.UserID)
				if err != nil {
					continue
				}
				if isTeamReply {
					message.Status = "responded"
					responseTime := (globals.ParseDuration(reply.Timestamp) - globals.ParseDuration(message.Timestamp)) / 60
					log.Debug("Response time is : ", responseTime)
					message.ResponseTime = time.Duration(responseTime)
					break
				}
			}
		}
		if len(message.Reactions) > 0 {
			for _, reaction := range message.Reactions {
				if reaction.Name == "heavy_check_mark" {
					message.Status = "fixed"
					break
				}
			}
		}
		message.RemindAt = ""

		err = a.ESClient.AddMessage(message, documentID)
		if err != nil {
			log.Error("Got an error while saving message", err)
			return []globals.SlackResponse{reply}, nil
		}
	}
	return []globals.SlackResponse{reply}, nil
}
