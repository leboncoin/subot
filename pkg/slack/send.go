package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// postAPIPayload posts an API request to Slack
func postAPIPayload(host string, endpoint string, payload string, channelToken string) error {
	queryURL := url.URL{Scheme: "https", Host: host, Path: fmt.Sprintf("api/%s", endpoint)}
	return postRequest(queryURL, payload, channelToken)
}

// postResponseURLPayload posts a request to Slack from a response URL
func postResponseURLPayload(responseURL string, payload string, channelToken string) error {
	queryURL, err := url.ParseRequestURI(responseURL)
	if err != nil {
		return err
	}
	return postRequest(*queryURL, payload, channelToken)
}

func postRequest(queryURL url.URL, payload string, channelToken string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", queryURL.String(), strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Bearer "+channelToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if resp.StatusCode != 200 {
		log.Errorf("%s %d", queryURL.String(), resp.StatusCode)
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error while reading response body")
		return err
	}

	var r ApiResponse
	if err = json.Unmarshal(bodyBytes, &r); err != nil {
		log.Error("An error occurred while parsing response body")
		return err
	}
	log.WithFields(log.Fields{"response": r}).Debug("reponse body from slack")
	if !r.Ok {
		return fmt.Errorf("Error while sending payload %s", r.Error)
	}
	return nil
}

// SendMessage calls Slack API on given channel URL with given body
func (s *Slack) SendMessage(text string, blocks []interface{}) error {
	payloadJSON := Event{
		Channel: s.Channel.ID,
		Text:    text,
		Blocks:  blocks,
	}
	payloadMarshalled, err := json.Marshal(payloadJSON)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)

	err = postAPIPayload(s.Host, "chat.postMessage", payloadString, s.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while posting api payload")
		return err
	}
	return nil
}

// ReplyToMessage calls Slack API on given channel URL with given body
func (s *Slack) ReplyToMessage(timestamp string, text string, blocks []interface{}) error {
	payloadJSON := Event{
		Blocks:   blocks,
		Channel:  s.Channel.ID,
		Text:     text,
		ThreadTs: timestamp,
		LinkNames: true,
	}

	payloadMarshalled, err := json.Marshal(payloadJSON)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)

	err = postAPIPayload(s.Host, "chat.postMessage", payloadString, s.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while posting api payload")
		return err
	}
	return nil
}

//DeleteResponseToMessage Calls the slack api to delete a message
func (s *Slack) DeleteResponseToMessage(timestamp string) error {
	payloadJSON := Event{
		Channel: s.Channel.ID,
		Ts:      timestamp,
	}
	payloadMarshalled, err := json.Marshal(payloadJSON)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)

	err = postAPIPayload(s.Host, "chat.delete", payloadString, s.BotToken)
	if err != nil {
		log.Fatal("Error while posting message to slack:", err)
	}
	return err
}

// SendEphemeralMessage sends an ephemeral message to the given user on the given channel
func (s *Slack) SendEphemeralMessage(userID string, text string) error {
	payloadJSON := Event{
		Channel: s.Channel.ID,
		User:    userID,
		Text:    text,
	}
	payloadMarshalled, err := json.Marshal(payloadJSON)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)

	err = postAPIPayload(s.Host, "chat.postEphemeral", payloadString, s.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while posting api payload")
		return err
	}
	return nil
}

// PostResponseURLPayload posts a request to Slack from a response URL
func (s *Slack) PostResponseURLPayload(responseURL string, text string) error {

	payload := UpdateBlockKit{
		ReplaceOriginal: true,
		Text:            text,
	}
	payloadMarshalled, err := json.Marshal(payload)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)
	err = postResponseURLPayload(responseURL, payloadString, s.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while posting response url payload")
		return err
	}
	return nil
}

// AddReaction places the requested emoji onto the message
func (s *Slack) AddReaction(timestamp string, name string) error {
	payloadJSON := Event{
		Channel:   s.Channel.ID,
		Name:      name,
		Timestamp: timestamp,
	}
	payloadMarshalled, err := json.Marshal(payloadJSON)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while marshalling json")
		return err
	}
	payloadString := string(payloadMarshalled)

	err = postAPIPayload(s.Host, "reactions.add", payloadString, s.BotToken)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Error while posting api payload")
		return err
	}
	return nil
}
