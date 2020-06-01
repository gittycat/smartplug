package smartplug

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	cmdInfo      string = `{"system":{"get_sysinfo":{}}}`
	cmdOn               = `{"system":{"set_relay_state":{"state":1}}}`
	cmdOff              = `{"system":{"set_relay_state":{"state":0}}}`
	cmdLedoff           = `{"system":{"set_led_off":{"off":1}}}`
	cmdLedon            = `{"system":{"set_led_off":{"off":0}}}`
	cmdCloudinfo        = `{"cnCloud":{"get_info":{}}}`
	cmdWlanscan         = `{"netif":{"get_scaninfo":{"refresh":0}}}`
	cmdTime             = `{"time":{"get_time":{}}}`
	cmdSchedule         = `{"schedule":{"get_rules":{}}}`
	cmdCountdown        = `{"count_down":{"get_rules":{}}}`
	cmdAntitheft        = `{"anti_theft":{"get_rules":{}}}`
	cmdReboot           = `{"system":{"reboot":{"delay":1}}}`
	cmdReset            = `{"system":{"reset":{"delay":1}}}`
	cmdEnergy           = `{"emeter":{"get_realtime":{}}}}`
)

// commands := map[string]string{
// 	"info":      `{"system":{"get_sysinfo":{}}}`,
// 	"on":        `{"system":{"set_relay_state":{"state":1}}}`,
// 	"off":       `{"system":{"set_relay_state":{"state":0}}}`,
// 	"ledoff":    `{"system":{"set_led_off":{"off":1}}}`,
// 	"ledon":     `{"system":{"set_led_off":{"off":0}}}`,
// 	"cloudinfo": `{"cnCloud":{"get_info":{}}}`,
// 	"wlanscan":  `{"netif":{"get_scaninfo":{"refresh":0}}}`,
// 	"time":      `{"time":{"get_time":{}}}`,
// 	"schedule":  `{"schedule":{"get_rules":{}}}`,
// 	"countdown": `{"count_down":{"get_rules":{}}}`,
// 	"antitheft": `{"anti_theft":{"get_rules":{}}}`,
// 	"reboot":    `{"system":{"reboot":{"delay":1}}}`,
// 	"reset":     `{"system":{"reset":{"delay":1}}}`,
// 	"energy":    `{"emeter":{"get_realtime":{}}}}`,
// }

func encrypt(str string) []byte {
	n := len(str)
	key := byte(171) // 171 is the hardcoded cypher in the smartplug. See doc
	result := make([]byte, n+4)
	// The length of the message is pre-pended to the encrypted payload
	// Also, the HS110 is using a MIPS chip (Big Endian based) so we
	// need to make sure that the 4 bytes (32 bit) number is encoded
	// correctly
	binary.BigEndian.PutUint32(result, uint32(n))
	for i, c := range str {
		result[i+4] = key ^ byte(c)
		key = result[i+4]
	}
	return result
}

func decrypt(buf []byte) string {
	key := byte(171)
	result := make([]byte, len(buf))
	for i, b := range buf {
		result[i] = key ^ b
		key = b
	}
	return string(result)
}

func send(addr string, cmd string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return "", fmt.Errorf("Could not connect to smartplug: %w", err)
	}
	defer conn.Close()

	request := encrypt(cmd)
	_, err = conn.Write(request)
	if err != nil {
		return "", fmt.Errorf("Could not write to smartplug: %w", err)
	}

	buf := make([]byte, 2048)
	_, err = conn.Read(buf) // Plug returns much less than 2048 bytes
	if err != nil {
		return "", fmt.Errorf("Could not read from smartplug: %w", err)
	}

	msg := decrypt(buf[4:]) // Strip first 4 bytes (length of message)
	return msg, nil
}
