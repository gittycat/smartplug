package smartplug

import (
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	in := `ABC abc 123 !@#$%"{}`
	enc := encrypt(in)
	out := decrypt(enc[4:])

	if in != out {
		t.Errorf("expected %s, got %s", in, out)
	}
}

func TestCmdInfo(t *testing.T) {
	msg, err := send("192.168.1.9:9999", cmdInfo)
	if err != nil {
		t.Errorf("Error sending %w", err)
	}

	fmt.Printf("Message: %s\n", msg)
}
