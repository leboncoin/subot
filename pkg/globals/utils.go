package globals

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

//DateLayout A constant for date formating
const DateLayout = "2006-01-02"

//CompareTimeStamps Utility to compare two timestamps and get the diff
func CompareTimeStamps(first string, second string) (float64, error) {
	firstInt, err := strconv.ParseFloat(first, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"timestamp": first,
			"error":     err,
		}).Error("Got an error while parsing first float")
		return 0., err
	}
	secondInt, err := strconv.ParseFloat(second, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"timestamp": second,
			"error":     err,
		}).Error("Got an error while parsing second float")
		return 0., err
	}
	return secondInt - firstInt, nil
}

//ParseDate converts a date to a timestamp
func ParseDate(date string) (string, error) {
	parsedDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", err
	}
	timestamp := strconv.FormatInt(parsedDate.Unix(), 10)
	return timestamp, nil
}

//ParseDuration parses string duration into float
func ParseDuration(t string) float64 {
	f, err := strconv.ParseFloat(t, 10)
	if err != nil {
		return 0.
	}
	return f
}
