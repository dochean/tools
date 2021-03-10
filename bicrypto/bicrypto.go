package bicrypto

import "errors"

type BiCryptor interface {
	// parse client public key
	SetClientPublicKey([]byte) error
	// Marshal Server Public key
	MarshalServerPublicKey() []byte
	// Encrypt client message
	EncryptClientMessage([]byte) []byte
	// Decrypt server message with private key
	DecryptServerMessage([]byte) []byte
}

type Bicrypto struct {
	s *RSACrypto
	c *RSACrypto
	Id string
	// expire time
}

func (b *Bicrypto) SetClientPublicKey(pkey []byte) error {
	b.c = NewRSAByPublicKey(pkey)
	if b.c == nil {
		return errors.New("Parse Client Public Key failed.")
	}
	return nil
}

func (b *Bicrypto) MarshalServerPublicKey() []byte {
	if b.s == nil {
		b.s = NewRSA()
	}
	return b.s.MarshalPublicKey()
}

func (b *Bicrypto) EncryptClientMessage(data []byte) []byte {
	if b.c == nil {
		return nil
	}
	return b.c.Encrypt(data)
}

func (b *Bicrypto) DecryptServerMessage(data []byte) []byte {
	if b.s == nil {
		return nil
	}
	return b.s.Decrypt(data)
}
