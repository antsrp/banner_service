package db

type Settings struct {
	Type     string `envconfig:"Type"`
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASS"`
	DBName   string `envconfig:"NAME"`
}

const (
	TypePostgres = "postgres"
	TypeMySql    = "mysql"
)
