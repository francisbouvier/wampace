package main

import (
	"log"

	"github.com/francisbouvier/wampace/client"
)

// Function to call when a published message is received
func event(args []interface{}, kwargs map[string]interface{}) {
	log.Println("Receive event with args:", args)
}

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

	// Subscribe to an event
	sub := c.Subscribe("com.myapp.mytopic1", event)

	// Check if subscription is done
	if err = <-sub.Subscribed; err != nil {
		log.Println("Subscription err")
	}

	// Acknowledgment are blocking:
	// use it inside a goroutine if you want them async, ie.
	// go func() {
	// 	if err = <-sub.Subscribed; err != nil {
	// 		log.Println("Subscription err")
	// 	}
	// }()

	// Unsubscribe
	// c.Unsubscribe(sub)
	// if err = <-sub.Unsubscribed; err != nil {
	// 	log.Println("Unsubscription error")
	// }

	// Wait for messages
	c.End()

}
