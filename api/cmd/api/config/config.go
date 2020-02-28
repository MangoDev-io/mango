package config

type Config struct {
	AlgodAddress string `required:"true"`
	PSToken      string `required:"true"`
	KMDAddress   string `required:"true"`
	KMDToken     string `required:"true"`

	DatabaseConfig
}

type DatabaseConfig struct {
	PostgreSQLUsername   string `required:"true"`
	PostgreSQLPassword   string `required:"true"`
	PostgreSQLDatabase   string `required:"true"`
	PostgreSQLHost       string `default:"db:5432"`
	MaxConnectionRetries int    `default:"15"`
}
