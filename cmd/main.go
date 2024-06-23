package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/WilliamTrojniak/TabAppBackend/cmd/api"
	"github.com/WilliamTrojniak/TabAppBackend/db"
	"github.com/WilliamTrojniak/TabAppBackend/env"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


func main() {

  databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", env.Envs.POSTGRES_USER, env.Envs.POSTGRES_PASSWORD, env.Envs.POSTGRES_HOST, env.Envs.POSTGRES_PORT, env.Envs.POSTGRES_DB);
  pgConfig, err := pgxpool.ParseConfig(databaseURL);
  if err != nil {
    log.Fatal(err);
  }

  pg, err := db.NewPostgresStorage(context.Background(), pgConfig);
  if err != nil {
    log.Fatal(err);
  }
  
  gob.Register(uuid.UUID{});
  server := api.NewAPIServer(":3000", pg.Pool);
  if err := server.Run(); err != nil {
    log.Fatal(err);
  }
  
}
