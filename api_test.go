package smartplug

import (
	"fmt"
	"testing"
	"time"
)

func TestEncryption(t *testing.T) {
	in := `ABC abc 0123 -+=:*!@#$%'"{}\n`
	enc := encrypt(in)
	out := decrypt(enc[4:])
	if in != out {
		t.Errorf("\"%s\"  !=  \"%s\"", in, out)
	}
}

func TestLEDs(t *testing.T) {
	p := NewSmartplug("192.168.1.9", "9999")

	err := p.LedOff()
	if err != nil {
		t.Errorf("Turning LED off. %w", err)
	}
	time.Sleep(2 * time.Second)
	err = p.LedOn()
	if err != nil {
		t.Errorf("Turning LED back on. %w", err)
	}
}

func TestGetInfo(t *testing.T) {
	p := NewSmartplug("192.168.1.9", "9999")

	info, err := p.Info()
	if err != nil {
		t.Errorf("requesting Info data. %w", err)
	}
	fmt.Printf("%v\n", info)
}

func TestGetMeter(t *testing.T) {
	p := NewSmartplug("192.168.1.9", "9999")

	info, err := p.Meter()
	if err != nil {
		t.Errorf("requesting Meter data. %w", err)
	}
	fmt.Printf("%+v\n", info)
}
