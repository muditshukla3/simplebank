package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/muditshukla3/simplebank/api"
	db "github.com/muditshukla3/simplebank/db/sqlc"
	"github.com/muditshukla3/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
		return
	}
	log.Println("config loaded")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server %v", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server %v", err)
	}
}
