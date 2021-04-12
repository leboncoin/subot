package elastic_test

import (
	"encoding/json"
	olivere "github.com/olivere/elastic"
	"net/http"
	"net/http/httptest"
	"github.com/leboncoin/subot/pkg/elastic"
	globals "github.com/leboncoin/subot/pkg/globals"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryTools(t *testing.T) {
	tool := "vault"
	expectedPath := "/tools/_search?pretty=true"

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

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.RequestURI, "Wrong path")
		res.WriteHeader(200)

		_, err := res.Write(expectedToolsJSON)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	hits, err := e.QueryTools(tool)
	assert.Equal(t, nil, err, "function shall not return errors")
	assert.Equal(t, 1, len(hits), "function shall not return no hits")
	assert.Equal(t, "mock", hits[0], "function shall return expected response")
}

func TestGetTools(t *testing.T) {
	expectedPath := "/tools/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took:     6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total:    1,
			MaxScore: 1.0,
			Hits: []elastic.Hit{{
				Index: "tools",
				Type:  "_doc",
				ID:    "mock0",
				Score: 1.0,
				Source: elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock0).*"}},
					},
				},
			}, {
				Index: "tools",
				Type:  "_doc",
				ID:    "mock1",
				Score: 1.0,
				Source: elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock1).*"}},
					},
				},
			}, {
				Index: "tools",
				Type:  "_doc",
				ID:    "mock2",
				Score: 1.0,
				Source: elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock2).*"}},
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
	hits, err := e.GetTools()
	assert.Equal(t, err, nil, "function shall not return errors")
	assert.Equal(t, len(hits), 3, "function shall not return no hits")
	assert.Equal(t, globals.Perco{
		ID:    "mock0",
		Query: globals.Query{Regexp: globals.Regexp{Input: ".*(mock0).*"}},
	}, hits[0], "function shall return expected response")
	assert.Equal(t, globals.Perco{
		ID:    "mock1",
		Query: globals.Query{Regexp: globals.Regexp{Input: ".*(mock1).*"}},
	}, hits[1], "function shall return expected response")
	assert.Equal(t, globals.Perco{
		ID:    "mock2",
		Query: globals.Query{Regexp: globals.Regexp{Input: ".*(mock2).*"}},
	}, hits[2], "function shall return expected response")
}

func TestAddTool(t *testing.T) {
	tool := globals.Perco{
		ID:    "mock",
		Name:  "mock",
		Query: globals.Query{
			Regexp: globals.Regexp{
				Input: ".*(mock|mock2).*",
			},
		},
	}

	expectedPath := "/tools/_doc/mock?refresh=true"

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
	err = e.AddTool(tool)
	assert.Equal(t, nil, err, "function shall not return errors")
}
