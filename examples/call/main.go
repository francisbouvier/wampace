package main

import (
	"log"

	"github.com/francisbouvier/wampace/client"
)

// Message to publish
var callMsg = []interface{}{"Hello World"}

func main() {

	// Set the address of the router
	addr := "127.0.0.1:1234"

	// Define a realm
	realm := "realm1"

	// Connect to the server
	c, err := client.New(addr)
	if err != nil {
		log.Fatalln("Connect error:", err)
	}

	// Join the wamp session
	err = c.Join(realm)
	if err != nil {
		log.Fatalln("Join error:", err)
	}

	// Call a remote procedure
	rc := c.Call("com.myapp.myprocedure1", callMsg, map[string]interface{}{})
	go func() {
		<-rc.Result
		log.Println("Results", rc.Args, rc.Kwargs)
	}()

	// Wait for messages
	c.End()

}
