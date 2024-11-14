package main

import (
	"Friends/cmd/config"
	"Friends/logg"
	"Friends/server"
	"Friends/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	lg := logg.New()
	lg.Info("start server")

	confg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		lg.Error("load config err", err)
	}

	psql, err := storage.New(lg,
		confg.App.Development.Database.Username,
		confg.App.Development.Database.Password,
		confg.App.Development.Database.Address,
		confg.App.Development.Database.NameDatabase,
	)
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

	httpServer := server.NewServer(lg, confg.App.Development.Server.HTTPPort, psql)
	if err := httpServer.Run(); err != nil {
		lg.Error("Server failed to start", err)
		return
	}
	lg.Info("Shutting down")
	//TODO реализовать gracefull shutdown

}
