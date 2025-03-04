package main

import (
	"github.com/AndreySirin/Friends/internal/servisAuth"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreySirin/Friends/internal/config"
	"github.com/AndreySirin/Friends/internal/logg"
	"github.com/AndreySirin/Friends/internal/server"
	"github.com/AndreySirin/Friends/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	lg := logg.New()
	lg.Info("start server")
	configPath, err := config.PathConfig()
	cfg, err := config.LoadConfig(lg, configPath)
	if err != nil {
		lg.Error("load config err", "error", err)
	}

	psql, err := storage.New(lg,
		cfg.App.Development.Database.Username,
		cfg.App.Development.Database.Password,
		cfg.App.Development.Database.Address,
		cfg.App.Development.Database.NameDatabase,
	)
	if err != nil || psql == nil {
		lg.Error("Failed to connect to database", "error", err)
		return
	}

	defer func() {
		if err = psql.Close(); err != nil {
			lg.Error("Failed to close",
				"error", err)
		}
	}()

	if err = psql.Migrate(migrate.Up); err != nil {
		lg.Error("Failed to migrate", "error", err)
		return
	}
	hash := servisAuth.NewBcryptHasher(10)
	secretSalt := []byte("my_secret_salt")
	au := servisAuth.NewAuth(psql, hash, secretSalt)
	if au == nil {
		lg.Error("Failed to initialize Auth service")
		return
	}
	httpServer := server.NewServer(lg, cfg.App.Development.Server.HTTPPort, psql, au)
	if httpServer == nil {
		lg.Error("Failed to initialize HTTP server")
		return
	}

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
