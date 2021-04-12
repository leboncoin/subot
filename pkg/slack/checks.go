package slack

import (
	"regexp"
	"github.com/leboncoin/subot/pkg/globals"
)

// IsValidToken checks if given token is the same as the config token
func (s Slack) IsValidToken(request EventRequest) bool {
	return request.Token == s.Token
}

// IsWatchedChannel checks if given channelId is the same as the one watched in the config
func (s Slack) IsWatchedChannel(event Event) bool {
	return event.Channel == s.Channel.ID || event.Item.Channel == s.Channel.ID
}

// IsThreadMessage checks if event is from a thread reply
func IsThreadMessage(event Event) bool {
	return event.ThreadTs != ""
}

// IsBotMessage checks if event is a message from bot
func (s Slack) IsBotMessage(event Event) bool {
	return event.SubType == "bot_message" || event.User == s.BotID
}

// IsChannelJoinMessage checks if event is a new member event
func IsChannelJoinMessage(event Event) bool {
	return event.SubType == "channel_join"
}

// IsChannelLeaveMessage checks if event is a leaving member event
func IsChannelLeaveMessage(event Event) bool {
	return event.SubType == "channel_leave"
}

// IsDeletedMessage checks if event is a leaving member event
func IsDeletedMessage(event Event) bool {
	if event.SubType == "message_deleted" {
		return true
	}
	if event.SubType == "message_changed" && event.Message.Text == "This message was deleted." {
		return true
	}
	return false
}

// IsUpdatedMessage checks if event is a leaving member event
func IsUpdatedMessage(event Event) bool {
	if event.SubType == "message_changed" && event.Message.Text != "This message was deleted." {
		return true
	}
	return false
}

// IsReaction checks if event is from a reaction added
func IsReaction(event Event) bool {
	return event.Type == "reaction_added" || event.Type == "reaction_removed"
}

// GetFiremanID returns the ID of the fireman found in the topic
func GetFiremanID(text string) string {
	topicPattern := regexp.MustCompile(`<@(.{11})> set the channel topic:.*<@(.{11})>`)
	if result := topicPattern.FindAllStringSubmatch(text, -1); len(result) > 0 {
		return result[0][2]
	}
	return ""
}

// IsTopicChange checks if event is a topic change containing a firemanID
func IsTopicChange(event Event) bool {
	if firemanID := GetFiremanID(event.Text); firemanID != "" {
		return true
	}
	return false
}

// GetMessageType checks event's content and returns its type
func (s Slack) GetMessageType(e Event) globals.MessageType {
	if IsThreadMessage(e) {
		return globals.Thread
	}

	if s.IsBotMessage(e) {
		return globals.BotMessage
	}

	if IsChannelJoinMessage(e) {
		return globals.Join
	}

	if IsChannelLeaveMessage(e) {
		return globals.Left
	}

	if IsReaction(e) {
		return globals.MessageReaction
	}

	if IsTopicChange(e) {
		return globals.TopicChange
	}

	if IsDeletedMessage(e) {
		return globals.DeletedMessage
	}

	if IsUpdatedMessage(e) {
		return globals.UpdatedMessage
	}

	return globals.NewMessage
}
