package main

import (
	"context"
	"log"
	"strconv"

	"github.com/WilliamTrojniak/TabAppBackend/cmd/api"
	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/env"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


func main() {

  port, err := strconv.Atoi(env.Envs.POSTGRES_PORT);
  if err != nil {
    log.Fatal("Could not convert POSTGRES_PORT to int");
  }

  _, err = db.NewPostgresStorage(context.Background(), &pgxpool.Config{
    ConnConfig: &pgx.ConnConfig{
      Config: pgconn.Config{
        Host: env.Envs.POSTGRES_HOST,
        Port: uint16(port),
        Database: env.Envs.POSTGRES_DB,
        User: env.Envs.POSTGRES_USER,
        Password: env.Envs.POSTGRES_PASSWORD,
      },
    },
  });
  if err != nil {
    log.Fatal(err);
  }
  
  server := api.NewAPIServer(":3000")
  if err := server.Run(); err != nil {
    log.Fatal(err);
  }
  
}
