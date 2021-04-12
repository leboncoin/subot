package slack

import (
	"strconv"
	"time"

	"github.com/leboncoin/subot/pkg/globals"
)

// GetEvent returns a Message struct populated from event data and more
func (s Slack) GetEvent(e Event) globals.Event {
	switch s.GetMessageType(e) {
	case globals.Thread:
		return s.GetReply(e)
	case globals.MessageReaction:
		return s.GetReaction(e)
	case globals.UpdatedMessage:
		return s.GetUpdatedMessage(e)
	default:
		return s.GetMessage(e)
	}
}

// GetMessage returns a Message struct populated from event data and more
func (s Slack) GetMessage(e Event) globals.Message {
	var apiResponse ApiResponse
	messageType := s.GetMessageType(e)
	if messageType == globals.TopicChange {
		firemanID := GetFiremanID(e.Text)
		apiResponse = s.ReadUser(firemanID)
	} else {
		apiResponse = s.ReadUser(e.User)
	}
	return globals.Message{
		Type:         s.GetMessageType(e),
		Text:         e.Text,
		UserID:       e.User,
		UserName:     apiResponse.User.Name,
		UserInfo:     apiResponse.User,
		Timestamp:    e.Ts,
		Reactions:    s.GetMessageReactions(e),
		Replies:      s.GetMessageReplies(e),
		ResponseTime: 0,
		DeletedTs:    e.DeletedTs,
		EditedTs:     e.Message.Timestamp,
	}
}

// GetUpdatedMessage returns a Message struct populated from update message Event (slightly different from message)
func (s Slack) GetUpdatedMessage(e Event) globals.Message {
	apiResponse := s.ReadUser(e.User)
	return globals.Message{
		Type:         s.GetMessageType(e),
		Text:         e.Message.Text,
		UserID:       e.Message.UserID,
		UserName:     apiResponse.User.Name,
		UserInfo:     apiResponse.User,
		Timestamp:    e.Message.Timestamp,
		Reactions:    s.GetMessageReactions(e),
		Replies:      s.GetMessageReplies(e),
		ResponseTime: 0,
		DeletedTs:    e.DeletedTs,
		EditedTs:     e.Message.Timestamp,
	}
}

// GetReaction returns a Reaction struct populated from event data and more
func (s Slack) GetReaction(e Event) globals.Reaction {
	return globals.Reaction{
		Name:      e.Reaction,
		MessageTs: e.Item.Ts,
		Users:     []string{e.User},
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Count:     0,
	}
}

// GetReply returns a Reply struct populated from event data and more
func (s Slack) GetReply(e Event) globals.Reply {
	apiResponse := s.ReadUser(e.User)
	return globals.Reply{
		ThreadTs:  e.ThreadTs,
		UserID:    e.User,
		UserName:  apiResponse.User.Name,
		UserInfo:  apiResponse.User,
		Timestamp: e.Ts,
		Text:      e.Text,
		FromBot:   e.User == s.BotID,
	}
}

// GetMessageReplies returns a list of Reply struct populated from all event data and more
func (s Slack) GetMessageReplies(e Event) (replies []globals.Reply) {
	if len(e.Message.Replies) > 0 {
		return e.Message.Replies
	}

	apiResponse, err := s.ReadReplies(e.Ts)

	if err != nil {
		return []globals.Reply{}
	}

	for _, message := range apiResponse.Messages {
		if message.ThreadTs == ""{
			continue
		}
		replies = append(replies, globals.Reply{
			ThreadTs:  message.ThreadTs,
			UserID:    message.User,
			UserName:  apiResponse.User.Name,
			UserInfo:  apiResponse.User,
			Timestamp: message.Ts,
			Text:      message.Text,
			FromBot:   message.User == s.BotID,
		})
	}
	return
}

// GetMessageReactions returns a list of Reactions struct populated from all event data and more
func (s Slack) GetMessageReactions(e Event) []globals.Reaction {
	if len(e.Message.Reactions) > 0 {
		return e.Message.Reactions
	}
	if len(e.Reactions) > 0 {
		return e.Reactions
	}
	return []globals.Reaction{}
}
