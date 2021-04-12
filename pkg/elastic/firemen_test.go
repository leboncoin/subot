package elastic_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/leboncoin/subot/pkg/elastic"
	"github.com/leboncoin/subot/pkg/globals"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryRangeFireman(t *testing.T) {
	start := "2020-01-01"
	end := "2020-01-07"
	expectedPath := "/firemen/_search?pretty=true"
	expectedMessage := globals.Message{
		Type:     "topic",
		Status:   "",
		Labels:   nil,
		Tools:    nil,
		Text:     "<@U3ZCU77BP> set the channel topic: <@U3ZCU77BP> est dédié au support cette semaine!\\n:redcard: Merci de ne pas utiliser de @here ou @channel\\n:point_right: Merci d'exposer ta question dans ton premier message\\n:threadplz: Merci de continuer la discussion en thread",
		UserID:   "UB210NGRK",
		UserName: "clement.mondion",
		UserInfo: globals.User{
			ID:         "UB210NGRK",
			SlackID:    0,
			Avatar:     "",
			TeamMember: false,
			Name:       "",
			Profile: globals.UserProfile{
				Email:                 "clement.mondion@adevinta.com",
				Avatar:                "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_32.jpg",
				Avatar512:             "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_512.jpg",
				RealName:              "Clément Mondion",
				LastName:              "Mondion",
				FirstName:             "Clément",
				DisplayName:           "clem",
				RealNameNormalized:    "Clement Mondion",
				DisplayNameNormalized: "clem",
			},
		},
		Timestamp:      "1592208201.000100",
		Reactions:      nil,
		Replies:        nil,
		EditedTs:       "",
		DeletedTs:      "",
		RemindAt:       "",
		ResponseTime:   0,
		ResolutionTime: 0,
		FeedbackStatus: "",
		FeedbackTs:     "",
	}

	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "firemen",
				Type:  "_doc",
				ID:    "ho4Bt3IBUCx0zk49He1Z",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: expectedMessage,
				},
			}},
		},
	}
	expectedJSONResponse, err := json.Marshal(expectedResponse)
	assert.Equal(t, nil, err, "Parsing json shall not return errors")

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.RequestURI, "Wrong path")
		res.WriteHeader(200)

		_, err := res.Write(expectedJSONResponse)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	hits, err := e.QueryRangeFireman(start, end)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 1, len(hits), "function shall not return no hits")
	assert.Equal(t, expectedMessage, hits[0], "function shall return expected response")
}

func TestAddFireman(t *testing.T) {
	message := globals.Message{
		Type:     "topic",
		Status:   "",
		Labels:   nil,
		Tools:    nil,
		Text:     "<@U3ZCU77BP> set the channel topic: <@U3ZCU77BP> est dédié au support cette semaine!\\n:redcard: Merci de ne pas utiliser de @here ou @channel\\n:point_right: Merci d'exposer ta question dans ton premier message\\n:threadplz: Merci de continuer la discussion en thread",
		UserID:   "UB210NGRK",
		UserName: "clement.mondion",
		UserInfo: globals.User{
			ID:         "UB210NGRK",
			SlackID:    0,
			Avatar:     "",
			TeamMember: false,
			Name:       "",
			Profile: globals.UserProfile{
				Email:                 "clement.mondion@adevinta.com",
				Avatar:                "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_32.jpg",
				Avatar512:             "https://avatars.slack-edge.com/2019-01-17/525636925408_32c9a24b43103a77c711_512.jpg",
				RealName:              "Clément Mondion",
				LastName:              "Mondion",
				FirstName:             "Clément",
				DisplayName:           "clem",
				RealNameNormalized:    "Clement Mondion",
				DisplayNameNormalized: "clem",
			},
		},
		Timestamp:      "1592208201.000100",
		Reactions:      nil,
		Replies:        nil,
		EditedTs:       "",
		DeletedTs:      "",
		RemindAt:       "",
		ResponseTime:   0,
		ResolutionTime: 0,
		FeedbackStatus: "",
		FeedbackTs:     "",
	}

	expectedPath := "/firemen/_doc/"

	serverResponse := map[string]interface{}{
		"_index": "firemen",
		"_type": "_doc",
		"_id": "I-LEfXQBBlaSKk1R5bDF",
		"_version": 1,
		"result": "created",
		"forced_refresh": true,
		"_seq_no": 4,
		"_primary_term": 2,
		"_shards": map[string]int{
			"total": 2,
			"successful": 1,
			"failed": 0,
		},
	}

	expectedJSONResponse, err := json.Marshal(serverResponse)
	assert.Equal(t, nil, err, "Parsing json shall not return errors")

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.RequestURI, "Wrong path")
		res.WriteHeader(201)

		_, err := res.Write(expectedJSONResponse)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	err = e.AddFireman(message)
	assert.Equal(t, nil, err, "function shall not return errors")
}
