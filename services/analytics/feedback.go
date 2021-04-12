package analytics

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

type feedbackTextSection struct {
	Type string            `json:"type"`
	Text map[string]string `json:"text"`
}

type feedbackElementSection struct {
	Type  string                     `json:"type"`
	Style string                     `json:"style"`
	Value string                     `json:"value"`
	Text  feedbackElementTextSection `json:"text"`
}

type feedbackElementTextSection struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

type feedbackActionsSection struct {
	Type     string                   `json:"type"`
	Elements []feedbackElementSection `json:"elements"`
}

func createFeedbackTemplate() []interface{} {
	template := []interface{}{
		feedbackTextSection{
			Type: "section",
			Text: map[string]string{
				"type": "mrkdwn",
				"text": "Cette réponse a t'elle permis de résoudre ton souci ?",
			},
		},
		feedbackActionsSection{
			Type: "actions",
			Elements: []feedbackElementSection{
				feedbackElementSection{
					Type:  "button",
					Style: "primary",
					Value: "feedback_useful",
					Text: feedbackElementTextSection{
						Type:  "plain_text",
						Text:  "Oui, merci Subot ! :slightly_smiling_face:",
						Emoji: true,
					},
				},
				feedbackElementSection{
					Type:  "button",
					Style: "danger",
					Value: "feedback_useless",
					Text: feedbackElementTextSection{
						Type:  "plain_text",
						Text:  "Non, contacter le pompier :fire:",
						Emoji: true,
					},
				},
			},
		},
	}
	return template
}

func getFeedbackResponse(ts string) *globals.SlackResponse {
	slackResponse := &globals.SlackResponse{
		Action: globals.ReplyMessage,
		Blocks: createFeedbackTemplate(),
		Ts:     ts,
	}
	return slackResponse
}


// handleFeedback godoc
// @Summary returns a reply when the interaction was completed
// @Description the reply returned is an update of the original feedback message
// @Description asking the user if the bot response was helpful.
// @Description Depending on the answer, it shall either call
// @Description the fireman for help or mark the original message as fixed.
// @Description Only the user that sent the original message
// @Description is allowed to complete the interaction
// @Tags Analytics
// @ID handle-feedback
// @Accept  json
// @Produce  json
// @Param interaction body object true "Interaction object sent by slack"
// @Router /analytics/feedback [post]
func (a Analyser) handleFeedback(interaction globals.Interaction) (replies []globals.SlackResponse, err error) {
	log.Debug("Get original message")
	originalMessages, err := a.ESClient.QueryRangeMessages(interaction.ThreadTs, interaction.ThreadTs)
	if len(originalMessages) == 0 {
		log.Debug("Original message not found, return")
		return
	}
	if originalMessages[0].UserID != interaction.ActionUserID {
		log.Error("Users don't match. Only the message owner can update its status")
		return replies, errors.New("users don't match. Only the message owner can update its status")
	}
	originalMessages[0].FeedbackStatus = globals.FeedbackStatus(interaction.ActionValue)
	originalMessages[0].FeedbackTs = interaction.ActionTs
	if err != nil {
		return replies, err
	}

	finalMessage := fmt.Sprintf("Je rends la main au pompier. Au secours <@%s>", a.getFiremanID())
	if globals.FeedbackStatus(interaction.ActionValue) == globals.UsefulFeedback {
		originalMessages[0].Status = "fixed"
		finalMessage = "Ravi d'avoir pu aider"
		// Set the heavy_check_mark
		replies = append(replies, globals.SlackResponse{
			Action: globals.React,
			Text:   "heavy_check_mark",
			Blocks: nil,
			Ts:     interaction.ThreadTs,
		})
	}
	// Update message content using response url
	replies = append(replies, globals.SlackResponse{
		Action:      globals.UpdateBlockKit,
		Text:        finalMessage,
		ResponseURL: interaction.ResponseURL,
	})

	err = a.ESClient.AddMessage(originalMessages[0], originalMessages[0].ID)
	if err != nil {
		log.Error("Got an error while saving message", err)
		return replies, nil
	}

	return replies, nil
}
