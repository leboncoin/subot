package analytics

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"strings"

	"github.com/leboncoin/subot/pkg/globals"
)

func buildReport(statistics globals.Statistics, pastStatistics globals.Statistics) (reportForm reportResponse) {
	responseTime := statistics.ResponseTime
	resolutionRate := statistics.ResolutionRate
	messagesDiff := len(statistics.Messages) - len(pastStatistics.Messages)
	resolutionDiff := statistics.ResolutionRate - pastStatistics.ResolutionRate
	messagesDiffLabel := "more"
	if messagesDiff < 0 {
		messagesDiff = int(math.Abs(float64(messagesDiff)))
		messagesDiffLabel = "less"
	}

	var firemen []string
	for _, f := range statistics.Firemen {
		firemen = append(firemen, fmt.Sprintf("<@%s>", f.ID))
	}
	reportForm = reportResponse{
		Blocks: []interface{}{
			reportTextSection{
				Type: "section",
				Text: map[string]string{
					"type": "plain_text",
					"text": "Here are the statistics of our performance on the support for the past week",
				},
			},
			reportFieldsSection{
				Type: "section",
				Fields: []map[string]string{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Messages:*\n%d messages this week", len(statistics.Messages)),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Resolution rate*\n%d%% fixed", resolutionRate),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Average response time*\n%d min", responseTime),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Messages evolution*\n%d messages %s compared to last week", messagesDiff, messagesDiffLabel),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Resolution rate compared to last week*\n%+d %%", resolutionDiff),
					},
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Firemen*\n%s", strings.Join(firemen, " ")),
					},
				},
			},
			reportTextSection{
				Type: "section",
				Text: map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("For more statistics see our *<%s|analytics dashboard>*", viper.GetString("front_url")),
				},
			},
		},
	}
	return
}
