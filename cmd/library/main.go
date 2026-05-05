package main

import "log"
import "rest_library/internal/server"
import "rest_library/internal/storage/in_memory"
import "rest_library/internal"

func main() {
	cfg := internal.ReadConfig()
	log.Printf("\nServer addr: %s\nServer port: %d\n\n", cfg.Addr, cfg.Port)
	log.Println("Library service started...")

	storage := in_memory.NewStorage()
	server := server.NewLibraryAPI(storage, cfg.Addr, cfg.Port)

	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Library service stopped...")
}
