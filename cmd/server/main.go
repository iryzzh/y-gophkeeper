package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iryzzh/y-gophkeeper/internal/config"
	"github.com/iryzzh/y-gophkeeper/internal/server"
	"github.com/iryzzh/y-gophkeeper/internal/services/item"
	"github.com/iryzzh/y-gophkeeper/internal/services/token"
	"github.com/iryzzh/y-gophkeeper/internal/services/user"
	"github.com/iryzzh/y-gophkeeper/internal/store"
	"github.com/iryzzh/y-gophkeeper/internal/store/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.NewServerConfig()
	if err != nil {
		panic(err)
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

	fmt.Println(cfg.Version.String())

	tokenSvc := token.NewService(
		st,
		cfg.Security.AtExpiresIn,
		cfg.Security.RtExpiresIn,
		[]byte(cfg.Security.AccessSecret),
		[]byte(cfg.Security.RefreshSecret),
	)

	userSvc := user.NewService(
		st,
		cfg.Security.HashMemory,
		cfg.Security.HashIterations,
		cfg.Security.HashParallelism,
		cfg.Security.SaltLength,
		cfg.Security.KeyLength,
	)

	itemSvc := item.NewService(st)

	srv := server.NewServer(&cfg.Web, tokenSvc, userSvc, itemSvc, true)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("server run: %v", err.Error())
	}

	return nil
}
