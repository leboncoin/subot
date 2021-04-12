package slack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (s Slack) curlAPI(urlPath string, query url.Values) ([]byte, error) {
	queryURL := url.URL{Scheme: "https", Host: s.Host, Path: urlPath, RawQuery: query.Encode()}
	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("Error while closing body %s", err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, err
}

//ReadMessages Call slack api to retrieve all the messages in a period
func (s Slack) ReadMessages(start string, end string, cursor string, limit ...int) (res ApiResponse, err error) {
	urlPath := "api/conversations.history"

	query := url.Values{}
	query.Set("token", s.BotToken)
	query.Set("channel", s.Channel.ID)
	query.Set("oldest", start)
	query.Set("latest", end)
	query.Set("cursor", cursor)
	if len(limit) > 0 {
		query.Set("limit", strconv.Itoa(limit[0]))
	}

	bodyBytes, err := s.curlAPI(urlPath, query)
	if err != nil {
		log.Error("An error occurred while performing the request", err)
		return
	}

	if err = json.Unmarshal(bodyBytes, &res); err != nil {
		log.Error("An error occurred while parsing response body", err)
		return
	}

	if !res.Ok {
		return res, errors.New(res.Error)
	}

	return
}

//ReadReplies Call the slack api to retrieve all the replies related to the message sent a the given timestamp
func (s Slack) ReadReplies(ts string) (res ApiResponse, err error) {
	urlPath := "api/conversations.replies"

	query := url.Values{}
	query.Set("token", s.BotToken)
	query.Set("channel", s.Channel.ID)
	query.Set("ts", ts)
	query.Set("limit", strconv.Itoa(100))

	bodyBytes, err := s.curlAPI(urlPath, query)
	if err != nil {
		log.Error("An error occurred while performing the request ", err)
		return
	}

	if err = json.Unmarshal(bodyBytes, &res); err != nil {
		log.Error("An error occurred while parsing response body", err)
		return
	}

	if !res.Ok {
		return res, errors.New(res.Error)
	}

	return
}

//ReadUser Call the slack api to retrieve all of the information related to the user
func (s Slack) ReadUser(id string) ApiResponse {
	urlPath := "api/users.info"

	query := url.Values{}
	query.Set("token", s.BotToken)
	query.Set("user", id)

	bodyBytes, err := s.curlAPI(urlPath, query)
	if err != nil {
		log.Error("An error occurred while performing the request", err)
	}

	var r ApiResponse
	if err = json.Unmarshal(bodyBytes, &r); err != nil {
		log.Error("An error occurred while parsing response body", err)
	}
	return r
}
