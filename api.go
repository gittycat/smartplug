// TP-Link Wi-Fi Smart Plug Protocol Client
// For use with TP-Link HS-100 or HS-110
//
// SPDX-License-Identifier: MIT-0

package smartplug

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Smartplug struct {
	address string // IP + Port  eg: 192.168.0.10:9999
}

type SmartplugInfo struct {
	SwVer        string `json:"sw_ver"`
	HwVer        string `json:"hw_ver"`
	TypeInternal string `json:"type"`
	Model        string `json:"model"`
	Mac          string `json:"mac"`
	DevName      string `json:"dev_name"`
	Alias        string `json:"alias"`
	RelayState   int    `json:"relay_state"`
	OnTime       int    `json:"on_time"`
	ActiveMode   string `json:"active_mode"`
	Feature      string `json:"feature"`
	Updating     int    `json:"updating"`
	IconHash     string `json:"icon_hash"`
	Rssi         int    `json:"rssi"`
	LedOff       int    `json:"led_off"`
	LongitudeI   int    `json:"longitude_i"`
	LatitudeI    int    `json:"latitude_i"`
	HwID         string `json:"hwId"`
	FwID         string `json:"fwId"`
	DeviceID     string `json:"deviceId"`
	OemID        string `json:"oemId"`
	NextAction   struct {
		TypeInternal int `json:"type"`
	} `json:"next_action"`
	ErrCode int `json:"err_code"`
}

type SmartplugMeter struct {
	VoltageMv int `json:"voltage_mv"`
	CurrentMa int `json:"current_ma"`
	PowerMw   int `json:"power_mw"`
	TotalWh   int `json:"total_wh"`
	ErrCode   int `json:"err_code"`
}

func NewSmartplug(ip string, port string) *Smartplug {
	// TODO: check whether ip and port are valid using net.IP.parse(str)
	addr := fmt.Sprintf("%s:%s", ip, port)
	return &Smartplug{addr, nil}
}

//
// Public API
// -------------------------------------------------------
// We only implement the main calls to extract info from the smartplug.
// The full list of calls is published at https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt

// Get System Info (Software & Hardware Versions, MAC, deviceID, hwID etc.)
func (p *Smartplug) Info() (*SmartplugInfo, error) {
	cmd := `{"system":{"get_sysinfo":{}}}`

	var info struct {
		System struct {
			GetSysInfo SmartplugInfo `json:"get_sysinfo"`
		} `json:"system"`
	}
	err := p.process(cmd, &info)
	if err != nil {
		return nil, err
	}
	return &info.System.GetSysInfo, nil
}

// Get Realtime Current and Voltage Reading
func (p *Smartplug) Meter() (*SmartplugMeter, error) {
	cmd := `{"emeter":{"get_realtime":{}}}`

	var info struct {
		EMeter struct {
			GetRealtime SmartplugMeter `json:"get_realtime"`
		} `json:"emeter"`
	}
	err := p.process(cmd, &info)
	if err != nil {
		return nil, err
	}
	return &info.EMeter.GetRealtime, nil
}

// Turn LED on
func (p *Smartplug) LedOn() error {
	cmd := `{"system":{"set_led_off":{"off":0}}}`
	_, err := send(p.address, cmd)
	return err
}

// Turn LED off
func (p *Smartplug) LedOff() error {
	cmd := `{"system":{"set_led_off":{"off":1}}}`
	_, err := send(p.address, cmd)
	return err
}

func (p *Smartplug) process(cmd string, result interface{}) error {
	ret, err := send(p.address, cmd)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(ret), result)
	if err != nil {
		fmt.Printf("can't unmarshal")
		return err
	}
	return nil
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

	buf := make([]byte, 2048) // No response is more than 2K long
	n, err := conn.Read(buf)  // Read up to 2K and return (do not wait for EOF)
	if err != nil {
		return "", fmt.Errorf("Could not read from smartplug: %w", err)
	}
	msg := decrypt(buf[4:n]) // Strip first 4 bytes (length of message)
	return msg, nil
}

//
// The smartplug uses a rudimentary encryption with hardcoded cypher (171).
// It expects requests to have the format:
//    [msg-size] [encrypted-json]
// where "msg-size" is a 4 byte unsigned integer (UInt32) corresponding to
// the length (in bytes) of the unencrypted request (the json payload).
// Endianicity must be taken into account for the 4 bytes integer since the
// smartplug uses a MIPS chip with Big Endian number packing whereas
// Intel and AMD chips use Little Endian.
// Endianicity is not an issue with the json payload. It uses UTF-8 with
// 2 bytes per character.
// Refer to the SoftsCheck link in the README.md for more details.
func encrypt(str string) []byte {
	n := len(str)
	key := byte(171) // cypher
	result := make([]byte, n+4)
	binary.BigEndian.PutUint32(result, uint32(n))
	for i, c := range str {
		result[i+4] = key ^ byte(c)
		key = result[i+4]
	}
	return result
}

// buf must not include the 4 bytes message length prefix
func decrypt(buf []byte) string {
	key := byte(171)
	result := make([]byte, len(buf))
	for i, b := range buf {
		// if b == 0x00 {
		// 	break
		// }
		result[i] = key ^ b
		key = b
	}
	return string(result)
}
