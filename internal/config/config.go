package config

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8081"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgres://user:password@localhost:5432/golang-diploma_db?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8080"`
}
