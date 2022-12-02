package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iryzzh/y-gophkeeper/internal/store"
	"github.com/iryzzh/y-gophkeeper/internal/store/sqlite"

	"github.com/iryzzh/y-gophkeeper/internal/client"
	"github.com/iryzzh/y-gophkeeper/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.NewClientConfig()
	if err != nil {
		return fmt.Errorf("config init failed: %v", err.Error())
	}

	var st store.Store
	switch cfg.DB.Type {
	case "sqlite3":
		st, err = sqlite.NewStore(cfg.DB.DSN, cfg.DB.MigrationsPath)
		if err != nil {
			return fmt.Errorf("store init: %v", err.Error())
		}
	default:
		return fmt.Errorf("not implemented DB type: %v", cfg.DB.Type)
	}

	app := client.NewClient(cfg, st)

	if err = app.Run(ctx); err != nil {
		return fmt.Errorf("run failed: %v", err.Error())
	}

	return nil
}
