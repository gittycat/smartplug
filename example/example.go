//
// Sample terminal client demonstrating how to use the API
//
package main

import (
	"fmt"
	"time"

	"github.com/gittycat/smartplug"
)

// Execute the callback function at intervals
func repeat(interval time.Duration, duration time.Duration, callback func()) {
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				callback()
			}
		}
	}()
	time.Sleep(duration)
	ticker.Stop()
	done <- true
}

func main() {

	p := smartplug.NewSmartplug("192.168.1.9", "9999") // ip, port

	// Get Device Information
	info, err := p.Info()
	if err != nil {
		fmt.Printf("error: %s", err)

	}
	fmt.Printf("Dev Name: %s\nVersion:  %s\n\n", info.Alias, info.SwVer)

	// Get Power consumption every 2 seconds
	fmt.Println("Use Control-C to quit")
	interval := 2 * time.Second
	duration := 2 * time.Minute
	repeat(interval, duration, func() {
		info, err := p.Meter()
		if err != nil {
			fmt.Printf("error retrieving voltage: %s\n", err)
			return
		}
		fmt.Printf("Power: %d mW  Current: %d mA\n", info.PowerMw, info.CurrentMa)
	})

}
