package config

type Config struct {
	AlgodAddress string `required:"true"`
	PSToken string `required:"true"`

	KMDAddress string `required:"true"`
	KMDToken string `required:"true"`
}