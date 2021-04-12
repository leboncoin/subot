package analytics

import (
	"github.com/leboncoin/subot/pkg/globals"
)

// HandleReportRequest godoc
// @Summary Generate a report for the given period
// @Description returns the report containing performance data
// @Description of the period and compared to last week
// @Tags Analytics
// @ID handle-report-request
// @Accept  json
// @Produce  json
// @Param start body string true "The start of the period to generate report for"
// @Param end body string true "The end of the period"
// @Router /analytics/report [post]
func (a Analyser) HandleReportRequest(start string, end string) (replies []globals.SlackResponse, err error) {
	var reply globals.SlackResponse
	reply.Action = globals.ChannelMessage
	statistics, err := a.Analyse(start, end)
	if err != nil {
		return nil, err
	}

	pastStart := subSevenDays(start)
	pastEnd := subSevenDays(end)

	pastStatistics, err := a.Analyse(pastStart, pastEnd)
	if err != nil {
		return nil, err
	}
	report := buildReport(statistics, pastStatistics)
	reply.Text = "Report"
	reply.Blocks = report.Blocks
	return []globals.SlackResponse{reply}, err
}
