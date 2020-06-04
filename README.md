# Go API to access the TP-Link HS110 SmartPlug
[![Go Report Card](https://goreportcard.com/badge/github.com/gittycat/smartplug)](https://goreportcard.com/report/github.com/gittycat/smartplug)


Simple, idiomatic and working API to access a TP-Link HS100 or HS110 smartplug. 

The focus is on presenting the minimum code needed to query the smartplug. There are about 
[50 calls documented](https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt) to control the device. Only a handful of calls are 
implemented here although extending it should be fairly trivial.

## How to Use

Here's a minimal program to retrieve Power and Current metrics.

```
package main

import "github.com/gittycat/smartplug"

public main() {
	p := smartplug.NewSmartplug("192.168.1.9", "9999") // ip, port
    info, err := p.Meter()
    if err != nil {
        fmt.Printf("error retrieving metrics: %s\n", err)
        return
    }
    fmt.Printf("Power: %d mW  Current: %d mA\n", info.PowerMw, info.CurrentMa)
}

```

The included example/example.go program shows how to retrieve the metrics at intervals.



## Credits
The encryption and decryption is based on the elegant python code from SoftsCheck
https://github.com/softScheck/tplink-smartplug

This API is possible due to the reverse engineering work done
by Lubomir Stroetmann and Tobias Esser at SoftsCheck.
https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/

They also include a command line python app to query and control the device.
