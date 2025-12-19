package config

type Config struct {
	Listen   string
	Postgres Postgres
	Redis    Redis
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}
type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var DefaultConfig Config = Config{

	Listen: "localhost:8080",
	Postgres: Postgres{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "pavan",
		Dbname:   "mydb",
	},

	Redis: Redis{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	},
}
