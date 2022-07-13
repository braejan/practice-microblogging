package entities

import "database/sql"

const (
	Host     string = "localhost"
	Port     int    = 54321
	DBName   string = "postgres"
	DBUser   string = "postgres"
	DBSecret string = "postgres"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres() (postgres *Postgres) {
	postgres = &Postgres{}
	return
}
