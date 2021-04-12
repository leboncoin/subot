package globals

import (
	pb "github.com/leboncoin/subot/pkg/engine_grpc_client/engine"
	"time"
)

// User represents the slack user that wrote a message
type User struct {
	ID         string      `json:"id"`
	SlackID    int         `json:"slack_id"`
	Avatar     string      `json:"avatar"`
	TeamMember bool        `json:"team_member"`
	Name       string      `json:"name"`
	Profile    UserProfile `json:"profile"`
}

// TeamMember represents the slack user that wrote a message
type TeamMember struct {
	ID         string      `json:"id"`
	SlackID    string         `json:"slack_id"`
	Name       string      `json:"name"`
}

// UserProfile contains profile information about the user
type UserProfile struct {
	Email                 string `json:"email"`
	Avatar                string `json:"image_32"`
	Avatar512             string `json:"image_512"`
	RealName              string `json:"real_name"`
	LastName              string `json:"last_name"`
	FirstName             string `json:"first_name"`
	DisplayName           string `json:"display_name"`
	RealNameNormalized    string `json:"real_name_normalized"`
	DisplayNameNormalized string `json:"display_name_normalized"`
}

// Statistics is the object containing all information about support in a period
type Statistics struct {
	ID             int           `json:"id"`
	Messages       []Message     `json:"messages"`
	ResponseTime   time.Duration `json:"response_time"`
	ResolutionTime time.Duration `json:"resolution_time"`
	ResolutionRate int           `json:"resolution_rate"`
	Firemen        []User        `json:"firemen"`
	Start          string        `json:"start"`
	End            string        `json:"end"`
}

// Message is the main structure representing a message
type Message struct {
	ID             string         `json:"id"`
	Type           MessageType    `json:"type"`
	Status         string         `json:"status"`
	Labels         []string       `json:"labels"`
	Tools          []string       `json:"tools"`
	AILabels       []pb.Category  `json:"ai_labels"`
	AITools        []pb.Category  `json:"ai_tools"`
	Text           string         `json:"text"`
	UserID         string         `json:"user"`
	UserName       string         `json:"user_name"`
	UserInfo       User           `json:"user_info"`
	Timestamp      string         `json:"ts"`
	Reactions      []Reaction     `json:"reactions"`
	Replies        []Reply        `json:"replies"`
	EditedTs       string         `json:"edited_ts"`
	DeletedTs      string         `json:"deleted_ts"`
	RemindAt       string         `json:"remind_at"`
	ResponseTime   time.Duration  `json:"response_time"`
	ResolutionTime time.Duration  `json:"resolution_time"`
	FeedbackStatus FeedbackStatus `json:"feedback_status"`
	FeedbackTs     string         `json:"feedback_ts"`
}

// Reaction is an icon placed on a message. All info comes from slack api except ts
type Reaction struct {
	Name      string   `json:"name"`
	Users     []string `json:"users"`
	MessageTs string   `json:"message_ts,omitempty"`
	Timestamp string   `json:"ts,omitempty"`
	Count     int      `json:"count"`
}

// Reply represents a message sent in a thread associated with a message
type Reply struct {
	UserID    string `json:"user"`
	UserName  string `json:"user_name"`
	UserInfo  User   `json:"user_info"`
	Timestamp string `json:"ts"`
	ThreadTs  string `json:"thread_ts"`
	Text      string `json:"text"`
	FromBot   bool   `json:"from_bot"`
}

// Answer represents a known answer matching a tool and a label
type Answer struct {
	ID       string `json:"id,omitempty"`
	Tool     string `json:"tool,omitempty"`
	Label    string `json:"label,omitempty"`
	Answer   string `json:"answer"`
	Feedback bool   `json:"feedback"`
}

// Perco represents a percolate query to match a label
type Perco struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Query Query  `json:"query"`
}

// Query represents a percolator query
type Query struct {
	Regexp Regexp `json:"regexp"`
}

// Regexp represents an input from the percolator
type Regexp struct {
	Input string `json:"input"`
}

// ResponseAction is used to set the action in SlackResponse
type ResponseAction string

const (
	// ChannelMessage Post a message in the main channel (e.g. reports)
	ChannelMessage ResponseAction = "channel"
	// ReplyMessage Reply to the original message, starting a thread
	ReplyMessage ResponseAction = "reply"
	// React Place an emoji onto a message to mark it as fixed (Not yet available)
	React ResponseAction = "react"
	// Ephemeral Send an ephemeral message dedicated to only one person, only visible to him, to remind him the rules
	Ephemeral ResponseAction = "ephemeral"
	// Nothing Do not do anything
	Nothing ResponseAction = "nothing"
	// DeleteMessage Delete one of your messages in a thread after the original message has been deleted
	DeleteMessage ResponseAction = "delete"
	// UpdateBlockKit Update a block kit in a message
	UpdateBlockKit ResponseAction = "update_block_kit"
)

// SlackResponse describes the data returned from analytics API
type SlackResponse struct {
	Action      ResponseAction `json:"action"`
	Text        string         `json:"text"`
	Blocks      []interface{}  `json:"blocks"`
	Ts          string         `json:"ts"`
	ChanID      string         `json:"chan_id"`
	UserID      string         `json:"user_id"`
	ResponseURL string         `json:"response_url"`
}

// MessageType represents the type of event we received
type MessageType string

const (
	//Thread message send in a thread - for the moment, we don't react to those
	Thread MessageType = "thread"
	//Join When someone joins the channel - send ephemeral message to remind the rules of the chan
	Join MessageType = "join"
	//Left When someone leaves the channel - do not do anything
	Left MessageType = "left"
	//BotMessage When a bot posted a message - we don't want to react to those
	BotMessage MessageType = "bot"
	//NewMessage When someone posts a new message in the main thread
	NewMessage MessageType = "user"
	//MessageReaction When someone adds an emoji to a message - consider message status change depending on the icon
	MessageReaction MessageType = "reaction"
	//TopicChange When one of the team member defines the fireman of the week on monday - Save fireman identity
	TopicChange MessageType = "topic"
	//DeletedMessage When someone deletes his message in the main thread - delete message posted by self to close the thread
	DeletedMessage MessageType = "deleted"
	//UpdatedMessage When someone updates its message - we need to take into account the change to update the info stored and analyse again the text
	UpdatedMessage MessageType = "updated"
)

// Event interface representing an incoming event
type Event interface {
	GetType() MessageType
	JSONData() []byte
}

// FeedbackStatus holds the current feedback status for a thread
type FeedbackStatus string

const (
	//NoFeedback no feedback has been asked
	NoFeedback FeedbackStatus = "no_feedback"
	//AskedFeedback the feedback has been asked but is not answered yet
	AskedFeedback FeedbackStatus = "asked_feedback"
	//UsefulFeedback the feedback is answered and has been useful
	UsefulFeedback FeedbackStatus = "feedback_useful"
	//UselessFeedback the feedback is answered and has been useless
	UselessFeedback FeedbackStatus = "feedback_useless"
)

// Interaction represents an interaction with a slack button from a user
type Interaction struct {
	MessageTs    string `json:"message_ts"`
	ThreadTs     string `json:"thread_ts"`
	ActionTs     string `json:"action_ts"`
	ActionUserID string `json:"action_user_id"`
	ActionValue  string `json:"action_value"`
	ResponseURL  string `json:"response_url"`
}
