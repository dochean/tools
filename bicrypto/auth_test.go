package bicrypto

import (
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	b := NewBiManager()
	b.Set(3 * time.Second)

	clientKey := `
-----BEGIN RSA Public Key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyMDN2IB/OH+5cdl+UdDt
t4yeEhW/fHdWXcfIjmQy/hNrFFi4Kz4UZcXWKT5Wv9lFJ+b0COtgKRI4E0Pvi/rm
HPSl0ZkDB4v3KSGElDBMgzGQSL4L7IeVaAqu9ZvKURPNa7TRpb4fnRQEc5V3xtLd
wznEdme5PehJiU7atbxOFVUKOfxpLVhMwDSlm0ZDCGse8IIx7OIzORFiufvhHb+Z
9YYps0lGw1RMt12oRHRyGP246pkLg39en2ywqngPhGTZH9Fqju2b/oNhdDylRIEH
iog6o7BojCFmncPSmn3p1SefJ3Xhbps81E/eR45rQVhKRxHbtEOPcCcEmvvLczvx
lwIDAQAB
-----END RSA Public Key-----`
	id, publickey := b.Add([]byte(clientKey))
	en := NewRSAByPublicKey(publickey)
	//en2 := NewRSAByPublicKey([]byte(clientKey))
	cipher := en.Encrypt([]byte("you nerver know"))

	res, _ := b.Decrypt(id, cipher)
	t.Log(string(res))
	ret, _ := b.Encrypt(id, []byte("to client"))
	t.Log(ret)

	time.Sleep(4*time.Second)
	b.Check()
	_, err := b.Decrypt(id, cipher)
	t.Log(err)
}