package slack_test

import (
	"github.com/leboncoin/subot/pkg/globals"
	"testing"

	"github.com/leboncoin/subot/pkg/slack"
)

func TestGetMessageType(t *testing.T) {
	join := slack.Event{
		SubType: "channel_join",
		Ts:      "12",
		User:    "EZHI35NN",
	}

	s := slack.Slack{
		BotID: "BotID",
	}

	joinResult := s.GetMessageType(join)
	if joinResult != globals.Join {
		t.Errorf("did not detect join message: %s is not equal to %s", joinResult, "join")
	}
	left := slack.Event{
		Ts:      "12",
		SubType: "channel_leave",
		User:    "EZHI35NN",
	}

	leftResult := s.GetMessageType(left)
	if leftResult != globals.Left {
		t.Errorf("did not detect left message: %s is not equal to %s", leftResult, "left")
	}

	bot := slack.Event{
		Ts:      "12",
		SubType: "bot_message",
		User:    "EZHI35NN",
	}

	botResult := s.GetMessageType(bot)
	if botResult != globals.BotMessage {
		t.Errorf("did not detect bot message: %s is not equal to %s", botResult, "bot")
	}

	subot := slack.Event{
		Ts:      "12",
		SubType: "message",
		User:    "BotID",
	}

	subotResult := s.GetMessageType(subot)
	if botResult != globals.BotMessage {
		t.Errorf("did not detect bot message: %s is not equal to %s", subotResult, "bot")
	}

	thread := slack.Event{
		Ts:       "12",
		Type:     "EZHI35NN",
		ThreadTs: "32345",
	}

	threadResult := s.GetMessageType(thread)
	if threadResult != globals.Thread {
		t.Errorf("did not detect thread message: %s is not equal to %s", threadResult, "thread")
	}

	message := slack.Event{
		Ts:   "12",
		User: "EZHI35NN",
	}

	messageResult := s.GetMessageType(message)
	if messageResult != globals.NewMessage {
		t.Errorf("did not detect message message: %s is not equal to %s", messageResult, "message")
	}
}
