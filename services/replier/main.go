package replier

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // blank import for remote
	"github.com/leboncoin/subot/pkg/config"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/leboncoin/subot/pkg/slack"
	"github.com/leboncoin/subot/pkg/vault"
)

// Run can be used to start the replier service from outside of the package
func Run() {
	config.Initialize()
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	// Configure vault client
	if viper.GetBool("vault_enabled") {
		_, err := vault.Configure()
		if err != nil {
			log.Fatal("Could not initialize vault client: ", err)
		}
	}

	// Init slack client
	s := slack.Slack{
		Host:     "slack.com",
		Channel:  slack.Chan{
			ID:      viper.GetString("slack_id"),
			Webhook: viper.GetString("slack_webhook"),
		},
		Token:    viper.GetString("slack_oauth_access_token"),
		BotToken: viper.GetString("slack_bot_user_oauth_access_token"),
		BotID:    viper.GetString("slack_bot_id"),
	}

	// Init replier
	replier := &Handler{Slack: &s, ApiUrl: viper.GetString("analytics_url")}

	runReportCron(replier)
	runReminderCron(replier)
	runAPI(replier)
}

func runReportCron(instance *Handler) {
	paris, _ := time.LoadLocation("Europe/Paris")
	c := cron.New(
		cron.WithLocation(paris),
	)

	if _, err := c.AddFunc("0 19 * * 5", instance.SendWeeklyReport); err != nil {
		log.Fatal("could not start cron job: ", err)
	}
	c.Start()
}

func runReminderCron(instance *Handler) {
	paris, _ := time.LoadLocation("Europe/Paris")
	c := cron.New(
		cron.WithLocation(paris),
	)

	if _, err := c.AddFunc("* 10-12,14-18 * * 1-5", instance.SendReminders); err != nil {
		log.Fatal("could not start cron job: ", err)
	}

	c.Start()
}
