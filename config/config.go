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
	Kavenegar struct {
		ApiKey string `env:"KAVENEGAR_API_KEY" env-default:"6B483932586E7562562F776479643931547767735275333257386D4B364C4B4457742F6368727A376E4A343D" env-description:"The API Key for Kavenegar SMS Provider"`
	}
}
