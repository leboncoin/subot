package elastic_test

import (
	"encoding/json"
	olivere "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/leboncoin/subot/pkg/elastic"
	"github.com/leboncoin/subot/pkg/globals"
	"testing"
)

func TestQueryAnswersVaultRights(t *testing.T) {
	tools := []string{"vault"}
	labels := []string{"rights"}
	expectedPath := "/answers/_search?pretty=true"
	expectedQuery := `{"from":0,"query":{"bool":{"filter":[{"terms":{"tool.keyword":["vault"]}},{"terms":{"label.keyword":["rights"]}}]}},"size":1000}`
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "answers",
				Type:  "_doc",
				ID:    "vault-rights",
				Score: 1.0,
				Source: elastic.HitSource{
					AnswerHitSource: elastic.AnswerHitSource{
						Tool:   "vault",
						Label:  "rights",
						Answer: "As-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
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
		body, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, nil, err, "Error in body decode")
		assert.Equal(t, expectedQuery, string(body), "Wrong body")
		res.WriteHeader(200)

		_, err = res.Write(expectedJSONResponse)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	results, err := e.QueryAnswers(tools, labels)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 1, len(results), "function shall not return no hits")
	assert.Equal(t, []globals.Answer{{
		Tool:   "vault",
		Label:  "rights",
		Answer: "As-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
	}}, results, "function shall return expected response")
}

func TestQueryAnswersHello(t *testing.T) {
	var tools []string
	labels := []string{"hello"}
	expectedPath := "/answers/_search?pretty=true"
	expectedQuery := `{"from":0,"query":{"bool":{"filter":{"terms":{"label.keyword":["hello"]}},"must_not":{"exists":{"field":"tool"}}}},"size":1000}`
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "answers",
				Type:  "_doc",
				ID:    "XEXXEF",
				Score: 1.0,
				Source: elastic.HitSource{
					AnswerHitSource: elastic.AnswerHitSource{
						Tool:   "",
						Label:  "hello",
						Answer: "Merci de nous exposer ton problème dans ton message",
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
		body, err := ioutil.ReadAll(req.Body)
		assert.Equal(t, nil, err, "Error in body decode")
		assert.Equal(t, expectedQuery, string(body), "Wrong body")
		res.WriteHeader(200)

		_, err = res.Write(expectedJSONResponse)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	results, err := e.QueryAnswers(tools, labels)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 1, len(results), "function shall not return no hits")
	assert.Equal(t, []globals.Answer{{
		Tool:   "",
		Label:  "hello",
		Answer: "Merci de nous exposer ton problème dans ton message",
	}}, results, "function shall return expected response")
}

func TestGetAnswers(t *testing.T) {
	expectedPath := "/answers/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "answers",
				Type:  "_doc",
				ID:    "vault-rights",
				Score: 1.0,
				Source: elastic.HitSource{
					AnswerHitSource: elastic.AnswerHitSource{
						Tool:   "vault",
						Label:  "rights",
						Answer: "As-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
					},
				},
			}, {
				Index: "answers",
				Type:  "_doc",
				ID:    "consul-rights",
				Score: 1.0,
				Source: elastic.HitSource{
					AnswerHitSource: elastic.AnswerHitSource{
						Tool:   "consul",
						Label:  "rights",
						Answer: "Besoin d'un nouveau token consul ? Regardes les reviews sur consul-config",
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
	hits, err := e.GetAnswers()
	assert.Equal(t, err, nil, "function shall not return errors")
	assert.Equal(t, len(hits), 2, "function shall not return no hits")
	assert.Equal(t, globals.Answer{
		ID:     "vault-rights",
		Tool:   "vault",
		Label:  "rights",
		Answer: "As-tu bien vérifié le path de ton secret ?\nFormat: `apps/team-<team name>/<app name>/<environment>/<secret name>`\n(sans `/` en début de path :wink:)\nPlus d information disponible dans cette <https://confluence.mpi-internal.com/display/LBCCORE/Vault|documentation>",
	}, hits[0], "function shall return expected response")
	assert.Equal(t, globals.Answer{
		ID:     "consul-rights",
		Tool:   "consul",
		Label:  "rights",
		Answer: "Besoin d'un nouveau token consul ? Regardes les reviews sur consul-config",
	}, hits[1], "function shall return expected response")
}

func TestAddAnswer(t *testing.T) {
	answer := globals.Answer{
		ID:       "I-LEfXQBBlaSKk1R5bDF",
		Tool:     "mockTool",
		Label:    "mockLabel",
		Answer:   "Mocking a label and a tool is great way to test this",
		Feedback: false,
	}
	h := json.RawMessage(`{"name": "mock"}`)

	expectedTools := olivere.SearchResult{
		Hits: &olivere.SearchHits{
			TotalHits: 1,
			Hits: []*olivere.SearchHit{{
				Id:     "Dont care",
				Source: &h,
			}},
		},
	}
	expectedToolsJSON, err := json.Marshal(expectedTools)
	assert.Equal(t, nil, err, "Parsing json shall not return errorsn err")

	expectedPaths := []string{"/tools/_search?pretty=true", "/labels/_search?pretty=true", "/answers/_doc/?refresh=true"}

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
		assert.Contains(t, expectedPaths, req.RequestURI, "Wrong path")
		if req.RequestURI != "" {
			res.WriteHeader(200)
			_, err = res.Write(expectedToolsJSON)
			assert.Equal(t, nil, err, "Parsing json shall not return errors")
		} else {
			res.WriteHeader(201)
			_, err = res.Write(expectedJSONResponse)
			assert.Equal(t, nil, err, "Parsing json shall not return errors")
		}

	}))

	e := MockClient(t, mockESServer)
	err = e.AddAnswer(answer)
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestEditAnswer(t *testing.T) {
	answer := globals.Answer{
		ID:       "I-LEfXQBBlaSKk1R5bDF",
		Tool:     "mockTool",
		Label:    "mockLabel",
		Answer:   "Mocking a label and a tool is great way to test this",
		Feedback: false,
	}

	expectedPaths := []string{"/tools/_search?pretty=true", "/labels/_search?pretty=true", "/answers/_doc/I-LEfXQBBlaSKk1R5bDF?refresh=true"}

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

	h := json.RawMessage(`{"name": "mock"}`)

	expectedTools := olivere.SearchResult{
		Hits: &olivere.SearchHits{
			TotalHits: 1,
			Hits: []*olivere.SearchHit{{
				Id:     "Dont care",
				Source: &h,
			}},
		},
	}
	expectedToolsJSON, err := json.Marshal(expectedTools)
	assert.Equal(t, nil, err, "Parsing json shall not return errorsn err")

	expectedJSONResponse, err := json.Marshal(serverResponse)
	assert.Equal(t, nil, err, "Parsing json shall not return errors")

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Contains(t, expectedPaths, req.RequestURI, "Wrong path")
		if req.RequestURI != "" {
			res.WriteHeader(200)
			_, err = res.Write(expectedToolsJSON)
			assert.Equal(t, nil, err, "Parsing json shall not return errors")
		} else {
			res.WriteHeader(201)
			_, err = res.Write(expectedJSONResponse)
			assert.Equal(t, nil, err, "Parsing json shall not return errors")
		}
	}))

	e := MockClient(t, mockESServer)
	err = e.EditAnswer("I-LEfXQBBlaSKk1R5bDF", answer)
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestDeleteAnswer(t *testing.T) {
	expectedPath := "/answers/_doc/mockTool-mockLabel?refresh=true"

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
	err = e.DeleteAnswer("mockTool-mockLabel")
	assert.Equal(t, nil, err, "function shall not return errors")
}
