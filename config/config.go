package config

type Config struct {
	Global struct {
		IsProductionMode bool `env:"IS_PRODUCTION_MODE" env-default:"false" env-description:"Is in production mode"`
	}
	COMonitoring struct {
		GinMode     string `env:"GIN_MODE" env-default:"debug" env-description:"Gin framework mode (release or debug)"`
		Host        string `env:"SERVER_HOST" env-default:"co_monitoring_server" env-description:"Host for CO monitoring"`
		ServingPort string `env:"SERVER_PORT" env-default:"8080" env-description:"Serving Port number for CO monitoring"`
		LogLevel    string `env:"LOG_LEVEL" env-default:"debug" env-description:"Log Level for application logger"`
	}
}
