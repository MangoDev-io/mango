package config

type Config struct {
	TestnetAlgodAddress string `required:"true"`
	MainnetAlgodAddress string `required:"true"`

	PSToken string `required:"true"`

	TokenAuthPassword string `required:"true"`
}
