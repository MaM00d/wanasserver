package main

import "math/rand"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
}

func register(firstName, lastName, username string) *User {
	return &User{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
	}
}
