package main

import (
	"log"
)

func main() {
	// fmt.Println("Yeah Buddy")
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.InitDb(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%+v\n", store)
	server := NewAPIServer(":3000", store)
	server.Run()
}
