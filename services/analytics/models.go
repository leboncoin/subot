package analytics

import (
	"github.com/leboncoin/subot/pkg/elastic"
	"github.com/leboncoin/subot/pkg/engine_grpc_client"
)

// APIRequest defines the structure of a call to the analytics api
type APIRequest struct {
	Label    string   `form:"label" json:"label"`
	Labels   []string `form:"labels" json:"labels"`
	Message  string   `form:"message" json:"message"`
	Regex    string   `form:"regex" json:"regex"`
	Tool     string   `form:"tool" json:"tool"`
	Tools    []string `form:"tools" json:"tools"`
	Status   string   `form:"status" json:"status"`
	User     string   `form:"user" json:"user"`
	ID       string   `form:"id" json:"id"`
	Answer   string   `form:"answer" json:"answer"`
	Feedback bool     `form:"feedback" json:"feedback"`
}

// Analyser is the main app struct
type Analyser struct {
	ESClient elastic.Interface          `json:"es_client"`
	Engine   engine_grpc_client.IEngine `json:"engine"`
}

type reportTextSection struct {
	Type string            `json:"type"`
	Text map[string]string `json:"text"`
}

type reportFieldsSection struct {
	Type   string              `json:"type"`
	Fields []map[string]string `json:"fields"`
}

type reportResponse struct {
	Blocks []interface{} `json:"blocks"`
}
