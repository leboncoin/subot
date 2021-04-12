package replier

import (
	log "github.com/sirupsen/logrus"
)

// SendReminders retrieves all reminders to send and then execute the slack action
func (h Handler) SendReminders() {
	log.Debug("Looking for reminders to send")
	reminders, err := h.callAnalyticsAPI("GET", "reminders", nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error occurred while calling reminders endpoint of the analytics service")
		return
	}
	log.WithFields(log.Fields{"reminders": reminders}).Debug("Got reminders to send")
	for _, reminder := range reminders {
		h.executeSlackAction(reminder)
	}
	return
}
