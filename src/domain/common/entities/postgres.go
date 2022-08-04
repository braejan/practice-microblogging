package entities

import "database/sql"

const (
	Host     string = "postgres-db.cdpc1ksupwh6.us-east-2.rds.amazonaws.com"
	Port     int    = 5432
	DBName   string = "postgresdb"
	DBUser   string = "postgres"
	DBSecret string = "postgres123"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres() (postgres *Postgres) {
	postgres = &Postgres{}
	return
}
