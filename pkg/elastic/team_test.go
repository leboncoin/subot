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

func TestIsTeamMember(t *testing.T) {
	userID := "mock"
	expectedPath := "/team/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "team",
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
	teamMember, err := e.IsTeamMember(userID)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, true, teamMember, "function shall return expected response")
}

func TestIsNotTeamMember(t *testing.T) {
	userID := "mock"
	expectedPath := "/team/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    0,
			MaxScore: 0,
			Hits:     []elastic.Hit{},
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
	teamMember, err := e.IsTeamMember(userID)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, false, teamMember, "function shall return expected response")
}

func TestAddTeamMember(t *testing.T) {
	user := globals.TeamMember{
		SlackID: "UB210NGRK",
		Name:    "clem",
	}

	expectedPath := "/team/_doc/?refresh=true"

	serverResponse := map[string]interface{}{
		"_index":         "firemen",
		"_type":          "_doc",
		"_id":            "I-LEfXQBBlaSKk1R5bDF",
		"_version":       1,
		"result":         "created",
		"forced_refresh": true,
		"_seq_no":        4,
		"_primary_term":  2,
		"_shards": map[string]int{
			"total":      2,
			"successful": 1,
			"failed":     0,
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
	err = e.AddTeamMember(user)
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestDeleteTeamMember(t *testing.T) {
	expectedPath := "/team/_doc/UB210NGRK?refresh=true"

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
	err = e.DeleteTeamMember("UB210NGRK")
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestGetTeamMembers(t *testing.T) {
	expectedPath := "/team/_search?pretty=true"

	message1 := globals.TeamMember{
		ID:       "chuz&fzofzo23R92I",
		SlackID:   "",
		Name:       "",
	}

	message2 := globals.TeamMember{
		ID:       "fzubzEZF49NFEziohzf",
		SlackID:   "",
		Name:       "",
	}

	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "team",
				Type:  "_doc",
				ID:    "chuz&fzofzo23R92I",
				Score: 1.0,
				Source: elastic.HitSource{},
			}, {
				Index: "messages",
				Type:  "_doc",
				ID:    "fzubzEZF49NFEziohzf",
				Score: 1.0,
				Source: elastic.HitSource{},
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
	hits, err := e.GetTeamMembers()
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 2, len(hits), "function shall not return no hits")
	assert.Equal(t, message1, hits[0], "function shall return expected response")
	assert.Equal(t, message2, hits[1], "function shall return expected response")
}
