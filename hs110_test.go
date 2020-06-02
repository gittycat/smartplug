package smartplug

import (
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

func TestCommands(t *testing.T) {

	skip := map[string]bool{
		"off":       true,
		"ledoff":    true,
		"antitheft": true,
		"reboot":    true,
		"reset":     true,
	}

	for name, cmd := range commands {

		// Skip the commands if it's in the exception list
		// We don't want to reset, reboot, etc...
		if skip[name] {
			continue
		}
		_, err := send("192.168.1.9:9999", cmd)
		if err != nil {
			t.Errorf("sending cmd %s  error: %w", name, err)
		}
	}
}
