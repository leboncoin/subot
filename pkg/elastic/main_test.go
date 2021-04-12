package elastic_test

import (
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/elastic"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	var mockESServer *httptest.Server
	mockESServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.RequestURI != "/v1/analytics/e" {
			res.WriteHeader(500)
		}
		res.WriteHeader(200)
	}))

	viper.Set("elastic_url", mockESServer.URL)

	_, err := elastic.Configure(true)
	if err != nil {
		log.Debug("Error")
	}
}

func MockClient(t *testing.T, testServer *httptest.Server) elastic.ES {
	viper.Set("elastic_url", testServer.URL)
	c, err := elastic.Configure(true)
	assert.Equal(t, nil, err, "initializing es client should not return errors")
	if err != nil {
		log.Debug("Error")
	}
	return c
}
