package elastic

import "github.com/leboncoin/subot/pkg/globals"

// Match the base object returned by the es api
type Match struct {
	Took     int     `json:"took"`
	TimedOut bool    `json:"timed_out"`
	Hits     HitList `json:"hits"`
}

// HitList a list of hits
type HitList struct {
	Total    int     `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

// Hit a single hit
type Hit struct {
	Index  string    `json:"_index"`
	Type   string    `json:"_type"`
	ID     string    `json:"_id"`
	Score  float64   `json:"_score"`
	Source HitSource `json:"_source"`
}

// HitSource the source of a hit. Can either be a message or a tool / label
type HitSource struct {
	globals.Message
	ToolHitSource
	AnswerHitSource
}

// ToolHitSource The query of the tool that matched (but we only care about the index in this case)
type ToolHitSource struct {
	Query HitQuery `json:"query"`
}

// HitQuery The regex of the query that matched percolator
type HitQuery struct {
	Regexp HitRegexp `json:"regexp"`
}

// HitRegexp The regex input that matched percolator
type HitRegexp struct {
	Input string `json:"input"`
}

// AnswerHitSource represents an answer object from ES
type AnswerHitSource struct {
	Tool   string `json:"tool"`
	Label  string `json:"label"`
	Answer string `json:"answer"`
}

// Interface set of wrapper around elastic package to ease mocking
type Interface interface {
	AddAnswer(globals.Answer) error
	AddFireman(globals.Message) error
	AddLabel(globals.Perco) error
	AddMessage(globals.Message, ...string) error
	AddTeamMember(globals.TeamMember) error
	AddTool(globals.Perco) error
	DeleteAnswer(string) error
	DeleteLabel(string) error
	DeleteMessage(string) error
	DeleteTeamMember(string) error
	DeleteTool(string) error
	EditAnswer(string, globals.Answer) error
	EditLabel(string, globals.Perco) error
	EditMessage(string, globals.Message) error
	EditTeamMember(string, globals.TeamMember) error
	EditTool(string, globals.Perco) error
	GetAnswers() ([]globals.Answer, error)
	GetLabels() ([]globals.Perco, error)
	GetTeamMembers() ([]globals.TeamMember, error)
	GetTools() ([]globals.Perco, error)
	IsTeamMember(string) (bool, error)
	QueryAnswers([]string, []string) ([]globals.Answer, error)
	QueryLabels(string) ([]string, error)
	QueryLabelByName(string) ([]globals.Perco, error)
	QueryLastUserMessages(string) ([]globals.Message, error)
	QueryRangeFireman(string, string) ([]globals.Message, error)
	QueryRangeMessages(string, string) ([]globals.Message, error)
	QueryReminderMessages() ([]globals.Message, error)
	QueryTools(string) ([]string, error)
	QueryToolByName(string) ([]globals.Perco, error)
}
