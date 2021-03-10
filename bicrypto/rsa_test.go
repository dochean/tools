package bicrypto

import (
	"testing"
)

func TestFlow(t *testing.T) {
	r := NewRSA()
	t.Log(string(r.MarshalPrivateKey()))
	data := []byte("testing is a good way to test your program.")
	out := r.Encrypt(data)
	t.Log(string(out))
	in := r.Decrypt(out)
	t.Log(string(in))
}

// better test case
