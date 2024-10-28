package main

import (
	"context"
	"log"
	"newsSite/internal/api"
	"newsSite/internal/pgx"
	"newsSite/internal/service"
	"newsSite/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	ctx := context.Background()
	pgxConnURL := "postgresql://postgres:admin@localhost:5432/todolist?sslmode=disable"
	pool, err := pgxpool.New(ctx, pgxConnURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	repo := pgx.New(pool)
	srv := service.NewService(repo)

	a := api.InitApi(srv)
	httpServer := server.NewServer(a.InitRouter())
	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

}
