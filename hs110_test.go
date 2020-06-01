package smartplug

import (
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
