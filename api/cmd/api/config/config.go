package config

type Config struct {
	TestnetAlgodAddress string `required:"true"`
	MainnetAlgodAddress string `required:"true"`

	PSToken string `required:"true"`

	TokenAuthPassword string `required:"true"`

	DatabaseConfig
}

type DatabaseConfig struct {
	PostgreSQLUsername   string `required:"true"`
	PostgreSQLPassword   string `required:"true"`
	PostgreSQLDatabase   string `required:"true"`
	PostgreSQLHost       string `default:"db:5432"`
	MaxConnectionRetries int    `default:"15"`
}
