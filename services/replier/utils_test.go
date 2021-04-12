package replier

import (
	"testing"

	"github.com/leboncoin/subot/pkg/slack"
)

func TestIsAuthorizedEvent(t *testing.T) {
	channel := slack.Chan{
		ID:      "dummy",
		Webhook: "test",
	}

	s := slack.Slack{
		Host:     "",
		Channel:  channel,
		Token:    "token",
		BotToken: "botToken",
		BotID:    "botId",
	}

	h := &Handler{Slack: &s, ApiUrl: "http://analytics:8080"}
	result := h.isAuthorizedEvent(slack.EventRequest{Token: "fake token"})
	if result == true {
		t.Errorf("Wrong token was accepted")
	}

	result2 := h.isAuthorizedEvent(slack.EventRequest{Token: "token", Event: slack.Event{Channel: "dummy"}})
	if result2 != true {
		t.Errorf("Right token was not accepted")
	}
}
