package main

import (
	"cloth-mini-app/internal/config"
	"cloth-mini-app/internal/storage/postgresql"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

func main() {
	log.Println("config initializing...")
	config := config.MustLoad()

	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		log.Panicf("failed to init postgresql storage %v", err)
		fmt.Println(config)
		os.Exit(0)
	}

	log.Print("migration up")

	err = goose.Up(storage.DB, "./migrations")
	if err != nil {
		log.Fatal(err)
	}
}
