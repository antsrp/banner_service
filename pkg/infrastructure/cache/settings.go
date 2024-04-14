package cache

type Settings struct {
	Type           string `envconfig:"TYPE"`
	Host           string `envconfig:"HOST"`
	Port           int    `envconfig:"PORT"`
	Password       string `envconfig:"PASS"`
	DBName         int    `envconfig:"DB"`
	ExpirationTime int    `envconfig:"EXPIRATION_TIME"`
}
