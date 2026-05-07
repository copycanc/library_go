package main

import (
	"log"
	"rest_library/internal/storage/db_storage"
	"rest_library/internal/storage/in_memory"
)
import "rest_library/internal/server"
import "rest_library/internal"

func main() {
	cfg := internal.ReadConfig()
	log.Printf("\nServer addr: %s\nServer port: %d\n\n", cfg.Addr, cfg.Port)
	log.Println("Library service started...")

	var repo server.Storage

	db, err := db_storage.NewStorage()
	if err != nil {
		repo = in_memory.NewStorage()
	} else {
		repo = db
	}

	server := server.NewLibraryAPI(repo, cfg.Addr, cfg.Port)

	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Library service stopped...")
}
