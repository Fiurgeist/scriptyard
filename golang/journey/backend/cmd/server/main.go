package main

import (
	"context"
	"fiurgeist/journey/internal/cache"
	"fiurgeist/journey/internal/metrics"
	"fiurgeist/journey/internal/queue"
	"fiurgeist/journey/internal/server"
	"fiurgeist/journey/internal/store"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Starting server...")
	// Capture SIGINT to for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	queue := queue.NewQueue()
	defer queue.Close()

	s, err := store.NewStore(queue)
	if err != nil {
		log.Printf("Error creating store: %v\n", err)
	}
	defer s.Close()

	metrics := metrics.NewMetrics()
	defer metrics.Close()

	cache := cache.NewCache(metrics, queue)

	config := server.Config{
		Addr:        ":8083",
		PublicDir:   "../public",
		PublicJsDir: "../public/static/js",
	}
	srv := server.NewHTTPServer(config, metrics, cache)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err.Error() == "http: Server closed" {
				log.Println("Server shut down")
			} else {
				log.Printf("Error while ListenAndServe: %v\n", err)
			}
		}
	}()
	log.Println("Server is ready")

	// Wait for SIGINT
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}
}
