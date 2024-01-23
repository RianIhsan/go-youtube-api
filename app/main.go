package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	ytChannelID := os.Getenv("YOUTUBE_CHANNEL_ID")
	if ytApiKey == "" && ytChannelID == "" {
		log.Fatal("env var is required")
	}
	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey, ytChannelID))

	return mux
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		signint := make(chan os.Signal, 1)
		signal.Notify(signint, os.Interrupt)
		signal.Notify(signint, syscall.SIGTERM)
		<-signint

		log.Println("service interrupted, received")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("error shutting down http server: %v", err)
		}

		log.Println("service shutdown complete")

		close(idleConnsClosed)
	}()

	log.Println("service start on port :8080")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("service stop")
}
