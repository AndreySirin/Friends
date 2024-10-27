package main

import (
	"github.com/AndreySirin/Friends/internal/config"
	"github.com/AndreySirin/Friends/internal/logg"
	"github.com/AndreySirin/Friends/internal/server"
	"github.com/AndreySirin/Friends/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lg := logg.New()
	lg.Info("start server")

	cfg, err := config.LoadConfig(lg, "config/config.yaml")
	if err != nil {
		lg.Error("load config err", "error", err)
	}

	psql, err := storage.New(lg,
		cfg.App.Development.Database.Username,
		cfg.App.Development.Database.Password,
		cfg.App.Development.Database.Address,
		cfg.App.Development.Database.NameDatabase,
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
	if _, err = psql.MigriteUP(); err != nil {
		lg.Error("Failed to migrate", "error", err.Error())
		return
	}
	httpServer := server.NewServer(lg, cfg.App.Development.Server.HTTPPort, psql)

	go func() {

		if err = httpServer.Run(); err != nil {
			lg.Error("Server failed to start", "error", err)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	go func() {
		<-stop
		err = httpServer.ShutDown()
		if err != nil {
			lg.Error("Failed to shutdown gracefully", "error", err)
		}
		lg.Info("Server gracefully stopped", "error", err)
		close(done)
	}()
	<-done
}
