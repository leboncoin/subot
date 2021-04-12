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

func TestQueryLastUserMessages(t *testing.T) {
	userID := "mock"
	expectedPath := "/messages/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "messages",
				Type:  "_doc",
				ID:    "chuz&fzofzo23R92I",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: globals.Message{
						ID:       "chuz&fzofzo23R92I",
						Type:     "topic",
						Status:   "",
						Labels:   nil,
						Tools:    nil,
						Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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
					},
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
	hits, err := e.QueryLastUserMessages(userID)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 1, len(hits), "function shall not return no hits")
	assert.Equal(t, expectedResponse.Hits.Hits[0].Source.Message, hits[0], "function shall return expected response")
}

func TestQueryRangeMessages(t *testing.T) {
	start := "2020-01-01"
	end := "2020-01-07"
	expectedPath := "/messages/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "messages",
				Type:  "_doc",
				ID:    "chuz&fzofzo23R92I",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: globals.Message{
						ID:       "chuz&fzofzo23R92I",
						Type:     "topic",
						Status:   "",
						Labels:   nil,
						Tools:    nil,
						Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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
					},
				},
			}, {
				Index: "messages",
				Type:  "_doc",
				Score: 1.0,
				ID:    "fzubzEZF49NFEziohzf",
				Source: elastic.HitSource{
					Message: globals.Message{
						ID:       "fzubzEZF49NFEziohzf",
						Type:     "topic",
						Status:   "",
						Labels:   nil,
						Tools:    nil,
						Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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
					},
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
	hits, err := e.QueryRangeMessages(start, end)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 2, len(hits), "function shall not return no hits")
	assert.Equal(t, expectedResponse.Hits.Hits[0].Source.Message, hits[0], "function shall return expected response")
	assert.Equal(t, expectedResponse.Hits.Hits[1].Source.Message, hits[1], "function shall return expected response")
}

func TestQueryReminderMessages(t *testing.T) {

	message1 := globals.Message{
		ID:       "chuz&fzofzo23R92I",
		Type:     "topic",
		Status:   "",
		Labels:   nil,
		Tools:    nil,
		Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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
	message2 := globals.Message{
		ID:       "fzubzEZF49NFEziohzf",
		Type:     "topic",
		Status:   "",
		Labels:   nil,
		Tools:    nil,
		Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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

	expectedPath := "/messages/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "messages",
				Type:  "_doc",
				ID:    "chuz&fzofzo23R92I",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: message1,
				},
			}, {
				Index: "messages",
				Type:  "_doc",
				ID:    "fzubzEZF49NFEziohzf",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: message2,
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
	hits, err := e.QueryReminderMessages()
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 2, len(hits), "function shall not return no hits")
	assert.Equal(t, message1, hits[0], "function shall return expected response")
	assert.Equal(t, message2, hits[1], "function shall return expected response")
}

func TestAddMessage(t *testing.T) {
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

	expectedPath := "/messages/_doc/"

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.RequestURI, "Wrong path")
		res.WriteHeader(201)

		_, err := res.Write(expectedJSONResponse)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	err = e.AddMessage(message)
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestDeleteMessage(t *testing.T) {
	ts := "1592208201.000100"

	expectedPath := "/messages/_doc/_delete_by_query"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "messages",
				Type:  "_doc",
				ID:    "chuz&fzofzo23R92I",
				Score: 1.0,
				Source: elastic.HitSource{
					Message: globals.Message{
						Type:     "topic",
						Status:   "",
						Labels:   nil,
						Tools:    nil,
						Text:     "Bonjour, pouvez-vous me donner les droits en lecture sur vault au path apps/team-engprod/support-analytics/prod/slack svp ?",
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
					},
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
	err = e.DeleteMessage(ts)
	assert.Equal(t, nil, err, "function shall not return errors")
}
