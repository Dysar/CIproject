package main

import (
	"fmt"
	"github.com/phayes/hookserve/hookserve"
)

func main() {
	server := hookserve.NewServer()
	server.Port = 8888
	server.Path = "/postreceive"
	server.Secret = "supersecretcode"
	server.GoListenAndServe()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
		//downloads files marked by latest commit
		//run the test in that repo
		//check the test status
	}
}
