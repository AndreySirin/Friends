package main

import (
	"Friends/logg"
	"Friends/server"
	"Friends/storage"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	address  = "127.0.0.1:5432"
	username = "myuser"
	password = "mypassword"
	database = "mydatabase"

	httpPort = ":8080"
)

func main() {

	lg := logg.New()
	lg.Info("start server")

	psql, err := storage.New(lg, username, password, address, database)
	if err != nil {
		lg.Error("Failed to connect to database",
			"error", err)
		return
	}

	defer func() {
		if err = psql.Close(); err != nil {
			lg.Error("Failed to close",
				"error", err)
		}
	}()

	if err = psql.DummyMigration(context.Background()); err != nil {
		lg.Error("Failed to migrate",
			"error", err)
		return
	}

	httpServer := server.NewServer(lg, httpPort)
	if err := httpServer.Run(); err != nil {
		lg.Error("Server failed to start", err)
		return
	}
	lg.Info("Shutting down")

}
