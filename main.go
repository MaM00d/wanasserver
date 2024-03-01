package main

import (
	"log"

	"Server/user"
)

func main() {
	// fmt.Println("Yeah Buddy")
	store, err := user.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.InitDb(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)
	server := user.NewAPIServer(":3000", store)
	server.Run()
}
