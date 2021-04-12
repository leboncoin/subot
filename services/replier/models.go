package replier

import "github.com/leboncoin/subot/pkg/slack"

// Handler is the main app struct
type Handler struct {
	Slack  slack.Interface `json:"slack"`
	ApiUrl string          `json:"api_url"`
}
