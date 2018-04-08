package main

import (
	"log"
	"net/http"

	"github.com/scukonick/geddis/config"
	"github.com/scukonick/geddis/db"
	sw "github.com/scukonick/geddis/serveryyy/go"
	"github.com/spf13/viper"
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

	log.Fatal(http.ListenAndServe(cfg.ServerAPI.ListenAddr, router))
}
