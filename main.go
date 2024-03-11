package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/erlanggawulung/shopifyx/api"
	db "github.com/erlanggawulung/shopifyx/db/sqlc"
	"github.com/erlanggawulung/shopifyx/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}
	store := db.NewStore(conn)
	runGinServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server:", err)
	}
	err = server.Start("0.0.0.0:8000")
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
