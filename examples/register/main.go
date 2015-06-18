package main

import (
	"log"

	"github.com/francisbouvier/wampace/client"
)

// Function to register for remote call
func procedure(args []interface{}, kwargs map[string]interface{}) ([]interface{}, map[string]interface{}) {
	log.Println("Receive call with args:", args)
	argsReturn := append(args, "from remote client")
	kwargsReturn := map[string]interface{}{}
	return argsReturn, kwargsReturn
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

	// Register a function
	rp := c.Register("com.myapp.myprocedure1", procedure)

	// Check if registration is done
	if err = <-rp.Registred; err != nil {
		log.Println("Registration err")
	}

	// Acknowledgment are blocking:
	// use it inside a goroutine if you want them async, ie.
	// go func() {
	// 	if err = <-rp.Registred; err != nil {
	// 		log.Println("Registration err")
	// 	}
	// }()

	// Unregister
	// c.Unregister(rp)
	// if err = <-rp.Unregistred; err != nil {
	// 	log.Println("Unregistration error")
	// }

	// Wait for messages
	c.End()

}
