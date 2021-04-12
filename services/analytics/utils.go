package analytics

import (
	"time"

	"github.com/leboncoin/subot/pkg/globals"

	log "github.com/sirupsen/logrus"
)

const reminderInterval = 1 * time.Hour

func subSevenDays(date string) string {
	parsedDate, err := time.Parse(globals.DateLayout, date)
	if err != nil {
		log.Error("Error while parsing date")
		return ""
	}
	sub := parsedDate.AddDate(0, 0, -7)
	return sub.Format(globals.DateLayout)
}

func (a Analyser) getFiremanID() string {
	startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format(globals.DateLayout)
	endOfWeek := time.Now().AddDate(0, 0, 1).Format(globals.DateLayout)
	fireman, err := a.ESClient.QueryRangeFireman(startOfWeek, endOfWeek)
	if err != nil {
		return ""
	}
	if len(fireman) < 1 {
		return ""
	}
	return fireman[0].UserInfo.ID
}
