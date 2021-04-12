package analytics

import (
	"github.com/gin-gonic/gin"
	pb "github.com/leboncoin/subot/pkg/engine_grpc_client/engine"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleMessage godoc
// @Summary Handles every user message sent in the main thread
// @Description Returns a reply depending on the message provided.
// @Description It shall not respond to team members, returning a "Nothing" action.
// @Description If the user already sent a message less than a minute ago
// @Description it will remind him to respect threads.
// @Description The message is then analysed, looking for known tools and labels.
// @Description If a known answer is found for those tools and labels,
// @Description it shall send this answer and ask a feedback from the user.
// @Description At the end of the analyse, the message is stored
// @Description in the database with all the information extracted.
// @Tags Analytics
// @ID handle-message
// @Accept  json
// @Produce  json
// @Param message body object true "Account ID"
// @Router /analytics/user [post]
func (a Analyser) HandleMessage(message globals.Message) ([]*globals.SlackResponse, error) {
	var reply globals.SlackResponse
	var documentID string
	replies := make([]*globals.SlackResponse, 0)
	replies = append(replies, &reply)
	reply.Action = globals.ReplyMessage
	reply.Ts = message.Timestamp
	reply.Text = ""
	message.FeedbackStatus = globals.NoFeedback
	log.Debug("Check if message comes from team member")
	isTeamMessage, err := a.ESClient.IsTeamMember(message.UserID)
	if err != nil {
		return replies, err
	}
	if isTeamMessage {
		log.Debug("Its a team message, set reply Action to Nothing")
		reply.Action = globals.Nothing
		message.Type = "team"
	}
	log.Debug("Check if message not respecting threads")
	LastUserMessages, err := a.ESClient.QueryLastUserMessages(message.UserID)
	log.Debug("Last user message", message.UserID)
	if err != nil {
		log.Error("Failed to fetch last user messages", err)
	}
	if len(LastUserMessages) > 0 {
		log.Debug("Consecutive message detected")
		log.Debug("Set reply text to : please respect threads")
		reply.Text = reply.Text + "Merci de respecter les threads."
	} else {
		reply.Text = reply.Text + "Merci pour ton message."
	}
	log.Debug("Get tools for event")
	//incidents, err := a.ESClient.QueryIncidents(tools)
	labels, err := a.ESClient.QueryLabels(message.Text)
	if err != nil {
		log.Error("Got an error while querying labels", err)
		return replies, nil
	}
	tools, err := a.ESClient.QueryTools(message.Text)
	if err != nil {
		log.Error("Got an error while querying tools", err)
		return replies, nil
	}

	log.WithFields(log.Fields{"tools": tools, "labels": labels}).Debug("Got tools and labels")

	answers, err := a.ESClient.QueryAnswers(tools, labels)
	if err != nil {
		log.Error("Got an error while querying answers ", err)
		return replies, nil
	}

	for _, a := range answers {
		log.Debug("Found predefined answer from elasticsearch")
		reply.Text = reply.Text + "\n" + a.Answer
		// Add a feedback request when a predefined answer is added to the message
		if a.Feedback {
			replies = append(replies, getFeedbackResponse(message.Timestamp))
			message.FeedbackStatus = globals.AskedFeedback
		}
	}

	aiTools, err := a.Engine.AnalyseMessageTools(&pb.Text{Text: message.Text})
	if err != nil {
		log.Error("Got an error while analysing message tools using AI", err)
		aiTools = []pb.Category{}
	}

	aiLabels, err := a.Engine.AnalyseMessageLabels(&pb.Text{Text: message.Text})
	if err != nil {
		log.Error("Got an error while analysing message labels", err)
		aiLabels = []pb.Category{}
	}

	log.WithFields(log.Fields{"tools": aiTools, "labels": aiLabels}).Debug("Got tools and labels from engines")

	message.AITools = aiTools
	message.AILabels = aiLabels
	message.Tools = tools
	message.Labels = labels
	message.Status = "unresponded"
	message.RemindAt = strconv.FormatInt(time.Now().Add(reminderInterval).Unix(), 10)

	log.Debug("Save or update message with document ID = ", documentID)
	err = a.ESClient.AddMessage(message, documentID)
	if err != nil {
		log.Error("Got an error while saving message", err)
		return replies, nil
	}
	return replies, nil
}


// EditMessage godoc
// @Summary Edit tools, labels or status of a message
// @Description For the given message timestamp,
// @Description store this information give in the payload
// @Tags Messages
// @ID edit-message
// @Produce  json
// @Param id body object true "Message ID"
// @Param type body object true "Message type (one of [user, team])"
// @Param status body object true "Message status [unresponded, responded, fixed, deleted]"
// @Param labels body object true "Labels extracted from the text or manually updated"
// @Param tools body object true "Tools extracted from the text or manually updated"
// @Param ai_labels body object true "Labels interpreted by the AI"
// @Param ai_tools body object true "Tools interpreted by the AI"
// @Param text body object true "Message content"
// @Param user_id body object true "ID of the user who posted the message"
// @Param user_name body object true "Name of the user who posted the message"
// @Param user_info body object true "Additional info about the user"
// @Param timestamp body object true "Timestamp at which the message was originally posted"
// @Param reactions body object true "List of reactions added to this message"
// @Param replies body object true "List of replies to this message"
// @Param edited_ts body object true "Timestamp at which the message was last edited"
// @Router /messages/:message_ts [put]
func (a Analyser) EditMessage(c *gin.Context) {
	var eventRequest globals.Message
	messageTs := c.Param("message_ts")
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.EditMessage(messageTs, eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Set the message status to deleted in the data storage
// @Tags Messages
// @ID delete-message
// @Produce  json
// @Router /messages/:message_ts [delete]
func (a Analyser) DeleteMessage(c *gin.Context) {
	messageTs := c.Param("message_ts")
	if err := a.ESClient.DeleteMessage(messageTs); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
