package main

import (
	"Friends/base"
	"Friends/logg"
	"Friends/server"
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	address  = "127.0.0.1:5432"
	username = "andrey"
	password = "mysecretpassword"
	database = "mydatabase"

	httpPort = ":8080"
)

func main() {

	logger := logg.New()
	logger.Info("start server")

	psql, err := base.NewBase(logger, username, password, address, database)
	if err != nil {
		logger.Error("Failed to connect to database",
			"error", err)
		return
	}

	defer func() {
		if err = psql.Close(); err != nil {
			logger.Error("Failed to close",
				"error", err)
		}
	}()

	if err = psql.DummyMigration(context.Background()); err != nil {
		logger.Error("Failed to migrate",
			"error", err)

		return
	}

	_, err = psql.Exec("INSERT INTO mens (name, description, price) "+
		"VALUES ($1, $2, $3)", "den", "kvn", 100)

	if err != nil {
		logger.Error("Failed to insert")
	}

	httpServer := server.NewServer(logger, httpPort)
	if err := httpServer.Run(); err != nil {
		logger.Error("Server failed to start", err)
		return
	}
	logger.Info("Shutting down")

}
