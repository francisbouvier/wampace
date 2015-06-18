package main

import (
	"log"

	"github.com/francisbouvier/wampace/client"
)

// Message to publish
var pubMsg = []interface{}{"Hello World"}

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

	// Publish an event
	pub := c.Publish("com.myapp.mytopic1", pubMsg, map[string]interface{}{})
	<-pub.Published

	// Acknowledgment are blocking:
	// use it inside a goroutine if you want them async, ie.
	go func() {
		if err = <-pub.Published; err != nil {
			log.Println("Publication err")
		}
	}()

	// Wait for messages
	c.End()

}
