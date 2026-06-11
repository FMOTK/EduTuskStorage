package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tuskstorage "httpServer/TuskStorage"
	"httpServer/myApi"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	storage := tuskstorage.NewTuskStorage(ctx)

	handler := myApi.NewHandler(ctx, storage)

	http.HandleFunc("/work", handler.CreateTuskHandler)
	http.HandleFunc("/work/status", handler.GetTuskStatusHandler)

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Printf("server started on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server start faled : %s", err.Error())
		}
	}()

	<-ctx.Done()
	log.Println("Server stoped by shutdown signal...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shatdown faled: %s", err.Error())
	}

	log.Println("Server stoped gracefully!")
}
