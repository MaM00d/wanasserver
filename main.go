package main

import (
	"log"

	"github.com/joho/godotenv"

	db "Server/eldb"
)

func main() {
	godotenv.Load()
	// fmt.Println("Yeah Buddy")
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	// if err := store.InitDb(); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%+v\n", store)
	server := NewApiServer(":3000", store)
	server.Run()
}
