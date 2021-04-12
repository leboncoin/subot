package replier

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
	"github.com/leboncoin/subot/pkg/slack"
)

// HandleNewEvent godoc
//// @Summary Parse slack request and calls corresponding function
//// @Description Given the information in the payload received from slack
//// @Description this function will determine the kind of request received
//// @Description and call the right analytics api endpoint
//// @ID handle-new-event
//// @Produce  json
//// @Param request query object true "The original slack request"
//// @Router /event [post]
func (h Handler) HandleNewEvent(request slack.EventRequest) {
	if !h.isAuthorizedEvent(request) {
		log.Debug("Not authorized request")
		return
	}

	event := h.Slack.GetEvent(request.Event)
	jsonBody := event.JSONData()
	res, err := h.callAnalyticsAPI("POST", string(event.GetType()), bytes.NewReader(jsonBody))
	if err != nil {
		log.Errorf("Error while fetching analytics api for %s endpoint: %s", string(event.GetType()), err)
	}
	log.WithFields(log.Fields{"res": res}).Debugf("Got results from analytics api %s endpoint", string(event.GetType()))
	for _, reply := range res {
		h.executeSlackAction(reply)
	}
	return
}

// HandleNewInteraction godoc
// @Summary Pass the interaction to the analytics api
// @ID handle-new-interaction
// @Produce  json
// @Param request query object true "The original slack request"
// @Router /interactivity [post]
func (h Handler) HandleNewInteraction(request slack.InteractivityRequest) {
	log.WithFields(log.Fields{"request": request}).Debug("Handle new interaction")

	// if !h.isAuthorizedEvent(request) {
	// 	log.Debug("Not authorized request")
	// 	return
	// }
	for _, action := range request.Actions {
		endpoint := "feedback"
		payload := globals.Interaction{
			MessageTs:    request.Message.Timestamp,
			ActionTs:     action.ActionTs,
			ActionUserID: request.User.ID,
			ActionValue:  action.Value,
			ThreadTs:     request.Message.ThreadTs,
			ResponseURL:  request.ResponseURL,
		}
		log.WithFields(log.Fields{"payload": payload}).Debug("Payload for analytics")
		jsonBody := payload.JSONData()
		res, err := h.callAnalyticsAPI("POST", endpoint, bytes.NewReader(jsonBody))
		if err != nil {
			log.Error("Error while fetching analytics interaction endpoint: ", err)
		}
		log.WithFields(log.Fields{"res": res}).Debug("Got results from analytics interaction endpoint")
		for _, reply := range res {
			h.executeSlackAction(reply)
		}
	}
	return
}

// executeSlackAction executes the action described in given response
func (h Handler) executeSlackAction(response globals.SlackResponse) {
	log.WithFields(log.Fields{"res": response}).Debug("Reply to message if necessary")
	if response.Action == globals.ChannelMessage {
		log.WithFields(log.Fields{"res": response}).Debug("Sending channel message")
		err := h.Slack.SendMessage(response.Text, response.Blocks)
		if err != nil {
			log.Error("Error while sending message to channel: ", err)
		}
		return
	}

	if response.Action == globals.ReplyMessage {
		err := h.Slack.ReplyToMessage(response.Ts, response.Text, response.Blocks)
		if err != nil {
			log.Error("Error while sending message to thread: ", err)
		}
		return
	}

	if response.Action == globals.Ephemeral {
		err := h.Slack.SendEphemeralMessage(response.UserID, response.Text)
		if err != nil {
			log.Error("Error while sending message: ", err)
		}
		return
	}

	if response.Action == globals.DeleteMessage {
		log.WithFields(log.Fields{"res": response}).Debug("Delete message reply")
		err := h.Slack.DeleteResponseToMessage(response.Ts)
		if err != nil {
			log.Error("Error while sending message to channel: ", err)
		}
		return
	}

	if response.Action == globals.UpdateBlockKit {
		err := h.Slack.PostResponseURLPayload(response.ResponseURL, response.Text)
		if err != nil {
			log.Error("Error while updating block kit: ", err)
		}
		return
	}

	if response.Action == globals.React {
		log.WithFields(log.Fields{"res": response}).Debug("Add reaction to message")
		err := h.Slack.AddReaction(response.Ts, response.Text)
		if err != nil {
			log.Error("Error while adding reaction to message: ", err)
		}
		return
	}
}

// HandleNewIncident is not implemented yet
func (h Handler) HandleNewIncident(request slack.CommandRequest) {
	log.WithFields(log.Fields{"request": request}).Debug("new incident")
}

// HandleNewRepo is not implemented yet
func (h Handler) HandleNewRepo(request slack.CommandRequest) {
	log.WithFields(log.Fields{"request": request}).Debug("new repository")
}
