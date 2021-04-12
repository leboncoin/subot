package analytics

import (
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleFiremanChange godoc
// @Summary Handles topic changes looking for a new fireman
// @Description returns an empty reply and stores the new fireman in storage
// @Tags Analytics
// @ID handle-fireman-change
// @Accept  json
// @Produce  json
// @Param message body object true "The new topic to extract information from"
// @Router /analytics/topic [post]
func (a Analyser) HandleFiremanChange(message globals.Message) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	reply.Action = globals.Nothing
	log.Debug("Check if message comes from team member")
	isTeamMessage, err := a.ESClient.IsTeamMember(message.UserID)
	if err != nil {
		return
	}
	if !isTeamMessage {
		log.Debug("Its not a team message, we cannot accept foreign firemen")
		return []globals.SlackResponse{reply}, nil
	}

	log.Debug("Save fireman")

	err = a.ESClient.AddFireman(message)
	if err != nil {
		log.Error("Got an error while saving fireman", err)
		return []globals.SlackResponse{reply}, nil
	}
	return []globals.SlackResponse{reply}, nil
}
