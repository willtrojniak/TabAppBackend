package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/willtrojniak/TabAppBackend/cmd/api"
	"github.com/willtrojniak/TabAppBackend/db"
	"github.com/willtrojniak/TabAppBackend/env"
	"github.com/willtrojniak/TabAppBackend/services/events"
	"github.com/willtrojniak/TabAppBackend/services/notifications"
)

var logLevels = map[string]slog.Level{
	"dev": slog.LevelDebug,
}

func main() {

	slog.SetLogLoggerLevel(logLevels[env.DEV])

	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", env.Envs.POSTGRES_USER, env.Envs.POSTGRES_PASSWORD, env.Envs.POSTGRES_HOST, env.Envs.POSTGRES_PORT, env.Envs.POSTGRES_DB)
	pgConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	pg, err := db.NewPostgresStorage(context.Background(), pgConfig)
	if err != nil {
		log.Fatal(err)
	}

	opts := redis.Options{
		Addr:     env.Envs.REDIS_ADDR,
		Password: "",
		DB:       0,
	}
	redis := redis.NewClient(&opts)

	eventDispatcher := events.NewEventDispatcher()

	notificationsService := notifications.NewNotificationService(slog.Default(), eventDispatcher)
	notificationsService.RegisterDriver(notifications.NewMailDriver(
		env.Envs.EMAIL_CLIENT_USER,
		env.Envs.EMAIL_CLIENT_PASSWORD,
		env.Envs.EMAIL_CLIENT_HOST,
		env.Envs.EMAIL_CLIENT_PORT),
		env.Envs.EMAIL_CLIENT_ENABLED)

	gob.Register(uuid.UUID{})
	server := api.NewAPIServer(":3000", pg, redis, eventDispatcher)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
