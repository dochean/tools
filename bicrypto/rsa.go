package bicrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var (
	KEY_LENGTH = 2048
	newHash    = sha256.New

	pubType = "RSA Public Key"
	prvType = "RSA Private Key"
)

type RSACrypto struct {
	prvKey *rsa.PrivateKey
	pubKey *rsa.PublicKey
}

func NewRSA() *RSACrypto {
	prvKey, err := rsa.GenerateKey(rand.Reader, KEY_LENGTH)
	check(err)
	return &RSACrypto{
		prvKey: prvKey,
		pubKey: &prvKey.PublicKey,
	}
}

func NewRSAByPrivateKey(prvKey []byte) *RSACrypto {
	prv := ParsePrivateKey(prvKey)
	if prv != nil {
		return &RSACrypto{
			prvKey: prv,
			pubKey: &prv.PublicKey,
		}
	}
	return nil
}

func NewRSAByPublicKey(pubKey []byte) *RSACrypto {
	pub := ParsePublicKey(pubKey)
	if pub != nil {
		return &RSACrypto{
			prvKey: nil,
			pubKey: pub,
		}
	}
	return nil
}

func (r *RSACrypto) MarshalPrivateKey() []byte {
	return marshalPrivateKey(r.prvKey)
}

func (r *RSACrypto) MarshalPublicKey() []byte {
	return marshalPublicKey(r.pubKey)
}

func (r *RSACrypto) Encrypt(data []byte) []byte {
	if r.pubKey == nil {
		fmt.Println("RSA Encryption Error: no public key")
		return nil
	}
	out, err := rsa.EncryptOAEP(newHash(), rand.Reader, r.pubKey, data, nil)
	if err != nil {
		fmt.Println("RSA Encryption Error: ", err)
		return nil
	}
	return out
}

func (r *RSACrypto) Decrypt(data []byte) []byte {
	if r.prvKey == nil {
		fmt.Println("RSA Decryption Error: no private key")
		return nil
	}
	out, err := rsa.DecryptOAEP(newHash(), rand.Reader, r.prvKey, data, nil)
	if err != nil {
		fmt.Println("RSA Decryption Error: ", err)
		return nil
	}
	return out
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func marshalPrivateKey(prvKey *rsa.PrivateKey) (prv []byte) {
	derStream := x509.MarshalPKCS1PrivateKey(prvKey)
	block := &pem.Block{
		Type:  prvType,
		Bytes: derStream,
	}
	prv = pem.EncodeToMemory(block)
	return
}

func marshalPublicKey(pubKey *rsa.PublicKey) (pub []byte) {
	derStream, err := x509.MarshalPKIXPublicKey(pubKey)
	check(err)
	block := &pem.Block{
		Type:  pubType,
		Bytes: derStream,
	}
	pub = pem.EncodeToMemory(block)
	return
}

func ParsePublicKey(pubKey []byte) (pub *rsa.PublicKey) {
	block, _ := pem.Decode(pubKey)
	if block == nil || block.Type != pubType {
		fmt.Println("Parse Public Key Err")
		return nil
	}
	pubItf, err := x509.ParsePKIXPublicKey(block.Bytes)
	check(err)
	pub = pubItf.(*rsa.PublicKey)
	return
}

func ParsePrivateKey(prvKey []byte) (prv *rsa.PrivateKey) {
	block, _ := pem.Decode(prvKey)
	if block == nil || block.Type != prvType {
		fmt.Println("Parse Private Key Err")
		return nil
	}
	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	check(err)
	return
}

func Encrypt(data []byte, pub []byte) (res []byte) {

	return
}

func Decrypt(data []byte, prv []byte) (res []byte) {

	return
}
