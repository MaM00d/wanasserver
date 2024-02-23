package main

// import "fmt"

func main() {
	// fmt.Println("Yeah Buddy")
	server := NewAPIServer(":3000")
	server.Run()
}
