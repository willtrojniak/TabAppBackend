package main

import (
	"log"
	"strconv"

	"github.com/WilliamTrojniak/TabAppBackend/cmd/api"
	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/env"
	"github.com/jackc/pgx"
)


func main() {

  config := env.GetConfig();
  port, err := strconv.Atoi(config.POSTGRES_PORT);
  if err != nil {
    log.Fatal("Could not convert POSTGRES_PORT to int");
  }

  db, err := db.NewPostgresStorage(pgx.ConnPoolConfig{
    ConnConfig: pgx.ConnConfig{
      Host: config.POSTGRES_HOST,
      Port: uint16(port),
      Database: config.POSTGRES_DB,
      User: config.POSTGRES_USER,
      Password: config.POSTGRES_PASSWORD,
    },
  });
  if err != nil {
    log.Fatal(err);
  }

  
  server := api.NewAPIServer(":3000", db)
  if err := server.Run(); err != nil {
    log.Fatal(err);
  }
  
}
