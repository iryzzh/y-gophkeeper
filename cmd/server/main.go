package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iryzzh/gophkeeper/internal/services/item"

	"github.com/iryzzh/gophkeeper/internal/services/user"

	"github.com/iryzzh/gophkeeper/internal/services/token"
	"github.com/iryzzh/gophkeeper/internal/store"
	"github.com/iryzzh/gophkeeper/internal/store/sqlite"

	"github.com/iryzzh/gophkeeper/internal/config"
	"github.com/iryzzh/gophkeeper/internal/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	var st store.Store
	switch cfg.DB.Type {
	case "sqlite3":
		st, err = sqlite.NewStore(cfg.DB.DSN, cfg.DB.MigrationsPath)
		if err != nil {
			panic("store: " + err.Error())
		}
	default:
		panic("init: not implemented DB type: " + cfg.DB.Type)
	}

	log.Println(cfg.GetVersion())

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

	srv := server.NewServer(&cfg.WebServer, tokenSvc, userSvc, itemSvc, true)

	var g errgroup.Group
	g.Go(func() error {
		return srv.Run(ctx)
	})

	if err = g.Wait(); err != nil {
		panic(err)
	}
}
