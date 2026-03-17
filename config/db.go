package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConectarBD() *sql.DB {
	// Substitua pela sua senha do Postgres
	stringConexao := "user=postgres dbname=app_login host=localhost port=5432 password=**** sslmode=disable"
	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	return db
}
