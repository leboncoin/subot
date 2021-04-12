package elastic_test

import (
	"encoding/json"
	olivere "github.com/olivere/elastic"
	"net/http"
	"net/http/httptest"
	"github.com/leboncoin/subot/pkg/elastic"
	"github.com/leboncoin/subot/pkg/globals"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryLabels(t *testing.T) {
	label := "rights"
	expectedPath := "/labels/_search?pretty=true"

	h := json.RawMessage(`{"name": "mock"}`)

	expectedLabels := olivere.SearchResult{
		Hits: &olivere.SearchHits{
			TotalHits: 1,
			Hits: []*olivere.SearchHit{{
				Id:     "Dont care",
				Source: &h,
			}},
		},
	}

	expectedLabelsJSON, err := json.Marshal(expectedLabels)
	assert.Equal(t, nil, err, "Parsing json shall not return errorsn err")

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.RequestURI, "Wrong path")
		res.WriteHeader(200)

		_, err := res.Write(expectedLabelsJSON)
		assert.Equal(t, nil, err, "Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	hits, err := e.QueryLabels(label)
	assert.Equal(t,  nil, err,"function shall not return errors")
	assert.Equal(t, 1, len(hits), "function shall not return no hits")
	assert.Equal(t, "mock", hits[0], "function shall return expected response")
}


func TestAddLabel(t *testing.T) {
	label := globals.Perco{
		ID:    "I-LEfXQBBlaSKk1R5bDF",
		Name:  "mock",
		Query: globals.Query{
			Regexp: globals.Regexp{
				Input: ".*(mock).*",
			},
		},
	}

	expectedPath := "/labels/_doc/I-LEfXQBBlaSKk1R5bDF?refresh=true"

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
	err = e.AddLabel(label)
	assert.Equal(t, nil, err, "function shall not return errors")
}

func TestGetLabels(t *testing.T) {
	expectedPath := "/labels/_search?pretty=true"
	expectedResponse := elastic.Match{
		Took: 6,
		TimedOut: false,
		Hits: elastic.HitList{
			Total : 1,
			MaxScore : 1.0,
			Hits : []elastic.Hit{{
				Index : "labels",
				Type : "_doc",
				ID : "mock0",
				Score : 1.0,
				Source : elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock0).*"}},
					},
				},
			},{
				Index : "labels",
				Type : "_doc",
				ID : "mock1",
				Score : 1.0,
				Source : elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock1).*"}},
					},
				},
			},{
				Index : "labels",
				Type : "_doc",
				ID : "mock2",
				Score : 1.0,
				Source : elastic.HitSource{
					ToolHitSource: elastic.ToolHitSource{
						Query: elastic.HitQuery{Regexp: elastic.HitRegexp{Input: ".*(mock2).*"}},
					},
				},
			}},
		},
	}
	expectedJSONResponse, err := json.Marshal(expectedResponse)
	assert.Equal(t, nil,err,  "Parsing json shall not return errors")

	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t,  expectedPath, req.RequestURI,"Wrong path")
		res.WriteHeader(200)

		_, err := res.Write(expectedJSONResponse)
		assert.Equal(t,  nil, err,"Parsing json shall not return errors")
	}))

	e := MockClient(t, mockESServer)
	hits, err := e.GetLabels()
	assert.Equal(t, err, nil, "function shall not return errors")
	assert.Equal(t, len(hits), 3, "function shall not return no hits")
	assert.Equal(t,  expectedResponse.Hits.Hits[0].ID, hits[0].ID, "function shall return expected response")
	assert.Equal(t, expectedResponse.Hits.Hits[1].ID, hits[1].ID, "function shall return expected response")
	assert.Equal(t,  expectedResponse.Hits.Hits[2].ID, hits[2].ID, "function shall return expected response")
}