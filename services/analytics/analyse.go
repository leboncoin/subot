package analytics

import (
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// Analyse godoc
// @Summary Retrieve analytics for the period
// @Description performs the analysis of the support performances
// @Description for the given period
// @Tags Analytics
// @ID analyse-messages
// @Accept  json
// @Produce  json
// @Param start query string true "Start date of the period to analyse (format 2020-12-31)"
// @Param end query string true "End date of the period to analyse (format 2020-12-31)"
// @Success 200 {object} map[string]string
// @Success 500 {string} Error
// @Router /analytics [get]
func (a Analyser) Analyse(start string, end string) (globals.Statistics, error) {
	messages, err := a.retrieveMessages(start, end)
	if err != nil {
		return globals.Statistics{}, err
	}
	firemen, err := a.retrieveFiremen(start, end)
	if err != nil {
		return globals.Statistics{}, err
	}
	responseTime, err := a.calculateResponseTime(messages)
	if err != nil {
		return globals.Statistics{}, err
	}
	resolutionRate, err := a.calculateResolutionRate(messages)
	if err != nil {
		return globals.Statistics{}, err
	}

	stats := globals.Statistics{
		Firemen:        firemen,
		Messages:       messages,
		ResponseTime:   responseTime,
		ResolutionRate: resolutionRate,
		Start:          start,
		End:            end,
	}
	return stats, nil
}

func (a Analyser) retrieveMessages(start string, end string) ([]globals.Message, error) {
	startTs, err := globals.ParseDate(start)
	if err != nil {
		return nil, err
	}
	endTs, err := globals.ParseDate(end)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"start": start,
		"end":   end,
	}).Debug("Starting fetch of messages to analyse")

	messages, err := a.ESClient.QueryRangeMessages(startTs, endTs)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (a Analyser) retrieveFiremen(start string, end string) ([]globals.User, error) {
	var firemen []globals.User
	esFiremen, err := a.ESClient.QueryRangeFireman(start, end)
	if err != nil {
		return nil, err
	}

	for _, fireman := range esFiremen {
		firemen = append(firemen, fireman.UserInfo)
	}

	return firemen, nil
}

func (a Analyser) calculateResponseTime(messages []globals.Message) (time.Duration, error) {
	log.Debug("Starting calculating response time")
	if len(messages) == 0 {
		return 0, nil
	}
	sum := 0.
	var responseTime []time.Duration
	for _, msg := range messages {
		if len(msg.Replies) == 0 {
			continue
		}
		responseTime = append(responseTime, msg.ResponseTime)
		sum = sum + float64(msg.ResponseTime)
	}
	log.Debug("Finished summing all individual response times")

	if sum == 0 {
		return 0, nil
	}

	log.Debug("Performing mean calculation")

	mean, err := time.ParseDuration(fmt.Sprintf("%f", sum/float64(len(responseTime))) + "ns")
	if err != nil {
		log.Errorf("got an error while calculating mean : %s", err)
		return 0, err
	}
	return mean, nil
}

func (a Analyser) calculateResolutionTime(messages []globals.Message) (time.Duration, error) {
	log.Debug("Starting calculating resolution time")
	if len(messages) == 0 {
		return 0, nil
	}
	sum := 0.
	var resolutionTime []time.Duration
	for _, msg := range messages {
		if len(msg.Reactions) == 0 {
			continue
		}
		resolutionTime = append(resolutionTime, msg.ResolutionTime)
		sum = sum + float64(msg.ResolutionTime)
	}
	log.Debug("Finished summing all individual response times")

	if sum == 0 {
		return 0, nil
	}

	log.Debug("Performing mean calculation")

	mean, err := time.ParseDuration(fmt.Sprintf("%f", sum/float64(len(resolutionTime))) + "ns")
	if err != nil {
		log.Errorf("got an error while calculating mean : %s", err)
		return 0, err
	}
	return mean, nil
}

func (a Analyser) calculateResolutionRate(messages []globals.Message) (int, error) {
	log.Debug("Starting calculating resolution rate")
	if len(messages) == 0 {
		return 0, nil
	}
	var fixedMessages []globals.Message
	log.Debug("Find proportion of fixed out of user message")
	for _, message := range messages {
		if message.Status == "fixed" {
			fixedMessages = append(fixedMessages, message)
		}
	}
	log.Debug("Finished categorising messages")
	if len(fixedMessages) < 1 {
		return 0, nil
	}
	return int(math.Min(100, float64(len(fixedMessages)*100/len(messages)))), nil
}
