package replier

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/globals"
)

// SendWeeklyReport builds the report and sends it
func (h Handler) SendWeeklyReport() {
	today := time.Now()
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday())+1).Format(globals.DateLayout)
	endOfWeek := today.AddDate(0, 0, 1).Format(globals.DateLayout)
	statsResponse, err := h.getStats(startOfWeek, endOfWeek)
	if err != nil {
		log.Fatal(err)
		return
	}
	h.executeSlackAction(statsResponse)
	return
}

func (h Handler) getStats(start string, end string) (response globals.SlackResponse, err error) {
	responses, err := h.callAnalyticsAPI("GET", fmt.Sprintf("report?start=%s&end=%s", start, end), nil)
	if err != nil {
		log.Error("Error while fetching analytics api for report endpoint: ", err)
	}
	log.WithFields(log.Fields{"res": response}).Debug("Got results from analytics api report endpoint")
	return responses[0], err
}
