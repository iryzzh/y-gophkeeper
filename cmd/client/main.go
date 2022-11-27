package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/iryzzh/gophkeeper/internal/client"
	"github.com/iryzzh/gophkeeper/internal/store/sqlite"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	st, err := sqlite.NewStore("gophkeeper.sqlite3", "../../migrations")
	if err != nil {
		panic("store: " + err.Error())
	}

	app := client.NewClient(ctx, st)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
