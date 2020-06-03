//
// Sample terminal client demonstrating how to use the API
//
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gittycat/smartplug"
)

func main() {

	p := smartplug.NewSmartplug("192.168.1.9", "9999")

	msg, err := p.Info()
	if err != nil {
		fmt.Printf("error: %s", err)

	}
	fmt.Printf("return = %s\n", msg)

	prettyJSON, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
