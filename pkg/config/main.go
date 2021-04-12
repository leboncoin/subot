package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Initialize configures the app
func Initialize() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Debug("Starting service")

	viper.SetDefault("env", "default")
	viper.AutomaticEnv()

	// Local configuration file
	viper.SetConfigName(viper.GetString("env"))      // name of config file (without extension)
	viper.AddConfigPath("config") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Errorf("Fatal error config file: %s", err)
	}
}
