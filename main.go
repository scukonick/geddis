package main

import (
	"log"
	"net/http"

	"os"
	"os/signal"

	"time"

	"github.com/scukonick/geddis/config"
	"github.com/scukonick/geddis/db"
	sw "github.com/scukonick/geddis/server_api/go"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func main() {
	log.Printf("Server started")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	cfg := &config.Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	store := db.NewGeddisStore(&cfg.Store)
	store.Run()

	router := sw.NewServerAPI(store).GetRouter()
	srv := &http.Server{
		Addr:         cfg.ServerAPI.ListenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("Received stop signal, exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	store.Stop()
	log.Println("shut down")
}
