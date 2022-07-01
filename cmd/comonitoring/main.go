package main

import (
	"comonitoring/config"
	"comonitoring/internal/comonitoring/api/rest"
	absms "comonitoring/internal/sms"
	"comonitoring/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

var log, _ = logger.NewLogger(logger.Config{})

func main() {
	log.Infof("Starting CO Monitoring Server...")

	// Read configs
	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)

	// Create the SMS provider
	kavenegar := absms.NewKavenegar(cfg.Kavenegar.ApiKey)

	// Create new Rest API server
	s, err := rest.NewApiServer(cfg, kavenegar)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Start serving
	log.Infof("The server is started on port %s.", cfg.COMonitoring.ServingPort)
	log.Fatalln(s.ListenAndServe().Error())
}
