package setup

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"log"
)

func ConnectToDB(env *Env) *sqlx.DB {
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	connection, err := pgx.ParseConfig(postgresURL)
	if err != nil {
		log.Fatalf("Error while parsing postgresUrl: %s", err)
	}
	conn, err := sqlx.Connect("pgx", connection.ConnString())
	if err != nil {
		log.Fatalf("Error while connecting to conn: %s", err)
	}

	return conn
}
