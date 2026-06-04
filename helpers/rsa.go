package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func GenerateRSAKey() (public, private string) {
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		fmt.Println("error generate")
	}

	// Extract public component.
	pub := key.Public()

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)
	// Encode private key to PKCS#1 ASN.1 PEM.
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	if er := ioutil.WriteFile("private"+".pem", privPEM, 0700); err != nil {
		println("error", er)
	}

	// Write public key to file.
	if errr := ioutil.WriteFile("public"+".pem", pubPEM, 0755); err != nil {
		println("error", errr)
	}
	return string(pubPEM), string(privPEM)
}

func ParsePublicKeyPem(pubPem string) (res *rsa.PublicKey) {
	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {
		fmt.Println("Failed")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse DER encoded public ")
	}

	return pub
}
func ParsePublicKey(pubPem string) (res *rsa.PublicKey) {
	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {
		fmt.Println("failed to parse PEM block containing the public key")
	}
	pkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	pub, ok := pkey.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("got unexpected key type: %T", pkey)
	}
	if err != nil {
		fmt.Println("failed to parse DER encoded public ")
	}

	return pub
}

func ParsePrivateKey(key string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		fmt.Println("failed to parse PEM block containing private key")
		return nil
	}

	if block.Type != "RSA PRIVATE KEY" {
		fmt.Println("unexpected private key type: ", block.Type)
		return nil
	}

	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("Error parse private key:", err.Error())
		return nil
	}

	//privateKey, ok := privkey.(*rsa.PrivateKey)
	//if !ok {
	//	log.Fatalf("got unexpected key type: %T", privkey)
	//}

	return privkey
}

func ParsePrivateKeyJws(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		fmt.Println("failed to parse PEM block containing private key")
		return nil, nil
	}

	fmt.Println(block.Type, "block =====")
	if block.Type != "PRIVATE KEY" {
		fmt.Println("unexpected private key type: ", block.Type)
		return nil, nil
	}

	privkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Println("Error parse private key:", err.Error())
		return nil, nil
	}

	//privateKey, ok := privkey.(*rsa.PrivateKey)
	//if !ok {
	//	log.Fatalf("got unexpected key type: %T", privkey)
	//}
	fmt.Println(privkey, "=== priv", reflect.TypeOf(privkey))
	if reflect.TypeOf(privkey).String() == "*rsa.PrivateKey" {
		return privkey.(*rsa.PrivateKey), nil
	}
	return nil, nil
}

func ParsePrivateKeyPem(privPEM string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	return priv
}
func GenerateSignatureSHA256(plaintext, privateKey string) (sig string) {
	log.Println("Generating signature SHA256...")
	msg := RsaMsg(plaintext)
	private := ParsePrivateKey(privateKey)
	signature, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, msg)
	if err != nil {
		fmt.Println("could not generate signature: ", err)
		return "-"
	}
	fmt.Println(msg, "--- sapi ---")
	fmt.Println(private, "--- private ---")
	return base64.StdEncoding.EncodeToString(signature)
}

func VerifSignatureSHA256(plaintext []byte, publicKey string, signature string) (is_ok bool) {
	pub := ParsePublicKey(publicKey)
	sig, _ := base64.StdEncoding.DecodeString(signature)
	//errs := rsa.VerifyPSS(pub, crypto.SHA256, plaintext, signature, nil)
	errs := rsa.VerifyPKCS1v15(pub, crypto.SHA256, plaintext, sig)
	if errs != nil {
		fmt.Println("could not verify signature: ", errs)
		return false
	}
	fmt.Println("signature verified")
	return true
}
func RsaMsg(text string) []byte {

	msg := []byte(text)
	msgHash := sha256.New()
	_, errs := msgHash.Write(msg)
	if errs != nil {
		return nil
	}

	msgHashSum := msgHash.Sum(nil)

	return msgHashSum
}
func RsaMsgSHA512(text string) []byte {

	msg := []byte(text)
	msgHash := sha512.New()
	_, errs := msgHash.Write(msg)
	if errs != nil {
		return nil
	}

	msgHashSum := msgHash.Sum(nil)

	return msgHashSum
}
func EncryptRSA(plaintext, publicKey string) (text []byte) {
	log.Println("Generating signature...")
	msg := RsaMsg(plaintext)
	pub := ParsePublicKeyPem(publicKey)
	//signature, err := rsa.SignPSS(rand.Reader, private, crypto.SHA256, msg, nil)
	text, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return
	}
	// text, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	// if err != nil {
	// 	return nil
	// }

	return
}
func DecryptRSA(plaintext, privateKey string) (text []byte) {
	log.Println("Generating signature...")
	msg := RsaMsg(plaintext)
	pub := ParsePrivateKeyPem(privateKey)
	//signature, err := rsa.SignPSS(rand.Reader, private, crypto.SHA256, msg, nil)
	text, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return
	}
	// text, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	// if err != nil {
	// 	return nil
	// }

	return
}
