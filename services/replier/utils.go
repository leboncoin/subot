package replier

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/leboncoin/subot/pkg/slack"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

//AnalyticsAPIResponse the structure of the responses of the analytics api
type AnalyticsAPIResponse struct {
	Error     string
	Responses []globals.SlackResponse
}

func (h Handler) isAuthorizedEvent(request slack.EventRequest) bool {
	// Verify token presence and consistency
	if !h.Slack.IsValidToken(request) {
		return false
	}

	// Check Channel ID
	if !h.Slack.IsWatchedChannel(request.Event) {
		log.Debug("Not from watched channel")
		return false
	}

	return true
}

func (h Handler) callAnalyticsAPI(method string, endpoint string, body io.Reader) (responses []globals.SlackResponse, err error) {
	url := h.ApiUrl + "/v1/analytics/" + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responses, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("Error while closing body %s", err)
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		err = fmt.Errorf(
			"invalid return code from analytics API : (%d) %s",
			resp.StatusCode,
			string(respBody),
		)
		return
	}

	err = json.Unmarshal(respBody, &responses)
	return
}
