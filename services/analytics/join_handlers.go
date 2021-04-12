package analytics

import (
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleJoinMessage godoc
// @Summary Handles people joining the channel
// @Description For newcomers, it returns an ephemeral message
// @Description containing the rules of the chan
// @Tags Analytics
// @ID handle-join-message
// @Accept  json
// @Produce  json
// @Param message body object true "Message content"
// @Router /analytics/join [post]
func (a Analyser) HandleJoinMessage(message globals.Message) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	reply.Action = globals.Ephemeral
	reply.UserID = message.UserID
	reply.Text = `:wave: _Bienvenue sur le chan de support engprod_

	:warning: *A lire avant de poster* :warning:

		:fireman: Une personne de l'équipe est dédiée chaque semaine à la gestion du support
		:redcard: Merci de ne pas utiliser de @here ou @channel
		:point_right: Merci d'exposer ta question dans ton premier message
		:threadplz: Merci de continuer la discussion en thread`

	return []globals.SlackResponse{reply}, nil
}
