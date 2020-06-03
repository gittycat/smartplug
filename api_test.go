package smartplug

import (
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	in := `ABC abc 0123 -+=:*!@#$%'"{}\n`
	enc := encrypt(in)
	out := decrypt(enc[4:])
	if in != out {
		t.Errorf("\"%s\"  !=  \"%s\"", in, out)
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
