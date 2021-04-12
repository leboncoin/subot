package slack

import "github.com/leboncoin/subot/pkg/globals"

// CommandRequest object received by slack when a user launches a slash command
type CommandRequest struct {
	Token       string `form:"token" json:"token"`
	Command     string `form:"command" json:"command"`
	ResponseURL string `form:"response_url" json:"response_url"`
	TeamID      string `form:"team_id" json:"team_id"`
	TeamDomain  string `form:"team_domain" json:"team_domain"`
	ChannelID   string `form:"channel_id" json:"channel_id"`
	ChannelName string `form:"channel_name" json:"channel_name"`
	UserID      string `form:"user_id" json:"user_id"`
	UserName    string `form:"user_name" json:"user_name"`
	Text        string `form:"text" json:"text"`
}

// Slack structure of the slack object
type Slack struct {
	Host     string `json:"host"`
	Channel  Chan   `json:"channel"`
	Token    string `json:"token"`
	BotToken string `json:"bot_token"`
	BotID    string `json:"bot_id"`
}

// Event wrapper around any kind of event received by slack
type Event struct {
	Ts              string             `form:"ts" json:"ts"`
	Timestamp       string             `form:"timestamp" json:"timestamp"`
	Type            string             `form:"type" json:"type"`
	SubType         string             `form:"subtype" json:"subtype"`
	Reaction        string             `form:"reaction" json:"reaction"`
	User            string             `form:"user" json:"user"`
	Text            string             `form:"text" json:"text"`
	Channel         string             `form:"channel" json:"channel"`
	EventTs         string             `form:"event_ts" json:"event_ts" json:"event_ts,omitempty"`
	ThreadTs        string             `form:"thread_ts" json:"thread_ts" json:"thread_ts,omitempty"`
	DeletedTs       string             `json:"deleted_ts" form:"deleted_ts" json:"deleted_ts,omitempty"`
	Blocks          []interface{}      `json:"blocks" json:"blocks,omitempty"`
	Message         globals.Message    `form:"message" json:"message"`
	Reactions       []globals.Reaction `form:"reactions" json:"reactions"`
	ReplyUsers      []string           `form:"reply_users" json:"reply_users"`
	PreviousMessage globals.Message    `form:"previous_message" json:"previous_message"`
	Item            EventItem          `form:"item" json:"item"`
	Name            string             `json:"name"`
	LinkNames       bool               `json:"link_names"`
}

// EventItem an Item of an even (used in some cases)
type EventItem struct {
	Type    string `form:"type" json:"type"`
	Channel string `form:"channel" json:"channel"`
	Ts      string `form:"ts" json:"ts"`
}

// EventRequest the request which wraps the Event
type EventRequest struct {
	Type      string `form:"type" json:"type"`
	Token     string `form:"token" json:"token"`
	Event     Event  `form:"event" json:"event"`
	Challenge string `form:"challenge" json:"challenge"`
}

// InteractivityChannel is the channel as represented in InteractivityRequest payload
type InteractivityChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// InteractivityAction is an action as represented in InteractivityRequest payload
type InteractivityAction struct {
	ActionID string `json:"action_id"`
	BlockID  string `json:"block_id"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	ActionTs string `json:"action_ts"`
}

// InteractivityRequest the request which wraps the payload from a slack interaction
type InteractivityRequest struct {
	Type        string                `json:"type"`
	User        globals.User          `json:"user"`
	APIAppID    string                `json:"api_app_id"`
	Token       string                `json:"token"`
	TriggerID   string                `json:"trigger_id"`
	Channel     InteractivityChannel  `json:"channel"`
	Message     globals.Reply         `json:"message"`
	ResponseURL string                `json:"response_url"`
	Actions     []InteractivityAction `json:"actions"`
}

// ResponseMetadata Metadata containing the cursor when fetching lots of data from the slack api
type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

// ApiResponse the response provided by the slack api when reading from it
type ApiResponse struct {
	Ok       bool             `json:"ok"`
	User     globals.User     `json:"user"`
	Error    string           `json:"error"`
	Latest   string           `json:"latest"`
	Oldest   string           `json:"oldest"`
	HasMore  bool             `json:"has_more"`
	Messages []Event          `json:"messages"`
	Metadata ResponseMetadata `json:"response_metadata"`
}

// Chan definition of a channel
type Chan struct {
	ID      string `json:"id"`
	Webhook string `json:"webhook"`
}

// Interface set of wrapper around slack package to ease mocking
type Interface interface {
	GetMessageType(e Event) globals.MessageType
	GetReply(e Event) globals.Reply
	GetMessage(e Event) globals.Message
	GetReaction(e Event) globals.Reaction
	GetUpdatedMessage(e Event) globals.Message
	GetEvent(e Event) globals.Event
	ReadUser(string) ApiResponse
	ReadMessages(string, string, string, ...int) (ApiResponse, error)
	SendMessage(string, []interface{}) error
	ReplyToMessage(string, string, []interface{}) error
	SendEphemeralMessage(string, string) error
	DeleteResponseToMessage(string) error
	IsValidToken(request EventRequest) bool
	IsWatchedChannel(event Event) bool
	PostResponseURLPayload(responseURL string, text string) error
	AddReaction(timestamp string, name string) error
}

// UpdateBlockKit represents the payload sent to a response url
type UpdateBlockKit struct {
	ReplaceOriginal bool   `json:"replace_original"`
	Text            string `json:"text"`
}
