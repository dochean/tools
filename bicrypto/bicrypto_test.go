package bicrypto

import "testing"

func TestBicrypto(t *testing.T) {
	bi := &Bicrypto{Id: "testing"}

	clientkey := `
-----BEGIN RSA Public Key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyMDN2IB/OH+5cdl+UdDt
t4yeEhW/fHdWXcfIjmQy/hNrFFi4Kz4UZcXWKT5Wv9lFJ+b0COtgKRI4E0Pvi/rm
HPSl0ZkDB4v3KSGElDBMgzGQSL4L7IeVaAqu9ZvKURPNa7TRpb4fnRQEc5V3xtLd
wznEdme5PehJiU7atbxOFVUKOfxpLVhMwDSlm0ZDCGse8IIx7OIzORFiufvhHb+Z
9YYps0lGw1RMt12oRHRyGP246pkLg39en2ywqngPhGTZH9Fqju2b/oNhdDylRIEH
iog6o7BojCFmncPSmn3p1SefJ3Xhbps81E/eR45rQVhKRxHbtEOPcCcEmvvLczvx
lwIDAQAB
-----END RSA Public Key-----`
	bi.SetClientPublicKey([]byte(clientkey))

	cipher := bi.EncryptClientMessage([]byte(`you never know`))
	t.Log(string(cipher))

	serverkey := bi.MarshalServerPublicKey()
	t.Log(string(serverkey))

	r := NewRSAByPublicKey(serverkey)
	ci := r.Encrypt([]byte(`you will never know`))

	ret := bi.DecryptServerMessage(ci)
	t.Log(string(ret))
}