package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
)

const ELECTOR1_MSP = "Org1MSP.eDUwOTo6Q049VXNlcjFAb3JnMS5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="
const ELECTOR2_MSP = "Org2MSP.eDUwOTo6Q049VXNlcjFAb3JnMi5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxPPW9yZzIuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="

// const ADMIN_PUB_KEY string = `-----BEGIN CERTIFICATE-----
// MIICKTCCAc+gAwIBAgIQTUFZM0uHkjSVzbXSGfdXAjAKBggqhkjOPQQDAjBzMQsw
// CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
// YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
// b3JnMS5leGFtcGxlLmNvbTAeFw0yMjA1MTUxMTI2MDBaFw0zMjA1MTIxMTI2MDBa
// MGsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
// YW4gRnJhbmNpc2NvMQ4wDAYDVQQLEwVhZG1pbjEfMB0GA1UEAwwWQWRtaW5Ab3Jn
// MS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABCVcwwpmCcxU
// oAdwCBJPr3kBaNPGpqFCzYXZ/zv0RNGeBm0Z07bA07lwhNZ6HtWIHjhRXbkKYM4i
// 49ctbnCxKS+jTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1Ud
// IwQkMCKAIJfLMIp/JAvDlsfhsnFqpClhPpt0IVJwFlZkSnh13wxcMAoGCCqGSM49
// BAMCA0gAMEUCIQC0cwSVvkx8oTh/87dERe7lnYDl5ZPyBuyBA5dSWs7s/AIgR889
// qwRxuMGZG6KsLXw4P9zdFccUKEIIweVuOMkO1J0=
// -----END CERTIFICATE-----`

// TODO: приватный ключ от захардкоженного выше. ПОка пускай тут полежит
// пока только для тестов
const ADMIN_PRIV_KEY string = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgKbFeJ8Dpu10b4+wA
2O+L1NWIMIVkJtbq6KKPAH1quyWhRANCAAQlXMMKZgnMVKAHcAgST695AWjTxqah
Qs2F2f879ETRngZtGdO2wNO5cITWeh7ViB44UV25CmDOIuPXLW5wsSkv
-----END PRIVATE KEY-----`

// func ExtractPubKeyFromCert(certPEM []byte) (*ecdsa.PublicKey, error) {
// 	block, _ := pem.Decode(certPEM)
// 	if block == nil || block.Type != "CERTIFICATE" {
// 		return nil, fmt.Errorf("failed to parse certificate PEM")
// 	}
// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cert.PublicKey.(*ecdsa.PublicKey), nil
// }

func GetAdminPub() *ecdsa.PublicKey {
	ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
	return &ecdsa.PublicKey{}
}

type Signer struct {
	priv *ecdsa.PrivateKey
}

func (s *Signer) SignElectorPayload(electionName, electorMSP string) ([]byte, error) {
	sig, err := ecdsa.SignASN1(rand.Reader, s.priv, HashElectorPayload(electionName, electorMSP))
	if err != nil {
		return nil, err
	}

	return SignatureToLowS(&s.priv.PublicKey, sig)
}

// TODO: пока только для теста
func SignTestMessage() {
	signer, err := NewSigner([]byte(ADMIN_PRIV_KEY))
	if err != nil {
		log.Panic(err)
	}

	signature, err := signer.SignElectorPayload("Best Crypto Currency", ELECTOR2_MSP)
	if err != nil {
		log.Panic(err)
	}

	pub, err := ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
	if err != nil {
		log.Panic(err)
	}

	payload := HashElectorPayload("Best Crypto Currency", ELECTOR2_MSP)

	ok := VerifySignature(pub, payload, signature)
	fmt.Println(ok)

	signatureHex := fmt.Sprintf("%x", signature)
	signatureHashBytes := sha256.Sum256([]byte(signatureHex))
	signatureHashHex := fmt.Sprintf("%x", signatureHashBytes[:])

	fmt.Println("signatureHashHex = ", signatureHashHex)
	fmt.Println("signatureHex = ", signatureHex)
}

func NewSigner(skPEM []byte) (*Signer, error) {
	block, _ := pem.Decode(skPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid PEM")
	}

	sk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	priv, ok := sk.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key")
	}

	return &Signer{
		priv: priv,
	}, nil
}

type ECDSASignature struct {
	R, S *big.Int
}

var (
	// curveHalfOrders contains the precomputed curve group orders halved.
	// It is used to ensure that signature' S value is lower or equal to the
	// curve group order halved. We accept only low-S signatures.
	// They are precomputed for efficiency reasons.
	curveHalfOrders = map[elliptic.Curve]*big.Int{
		elliptic.P224(): new(big.Int).Rsh(elliptic.P224().Params().N, 1),
		elliptic.P256(): new(big.Int).Rsh(elliptic.P256().Params().N, 1),
		elliptic.P384(): new(big.Int).Rsh(elliptic.P384().Params().N, 1),
		elliptic.P521(): new(big.Int).Rsh(elliptic.P521().Params().N, 1),
	}
)

func marshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ECDSASignature{r, s})
}

func unmarshalECDSASignature(raw []byte) (*big.Int, *big.Int, error) {
	// Unmarshal
	sig := new(ECDSASignature)
	_, err := asn1.Unmarshal(raw, sig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed unmashalling signature [%s]", err)
	}
	// Validate sig
	if sig.R == nil {
		return nil, nil, errors.New("invalid signature, R must be different from nil")
	}
	if sig.S == nil {
		return nil, nil, errors.New("invalid signature, S must be different from nil")
	}
	if sig.R.Sign() != 1 {
		return nil, nil, errors.New("invalid signature, R must be larger than zero")
	}
	if sig.S.Sign() != 1 {
		return nil, nil, errors.New("invalid signature, S must be larger than zero")
	}
	return sig.R, sig.S, nil
}

func SignatureToLowS(k *ecdsa.PublicKey, signature []byte) ([]byte, error) {
	r, s, err := unmarshalECDSASignature(signature)
	if err != nil {
		return nil, err
	}
	s, err = toLowS(k, s)
	if err != nil {
		return nil, err
	}
	return marshalECDSASignature(r, s)
}

// IsLow checks that s is a low-S
func isLowS(k *ecdsa.PublicKey, s *big.Int) (bool, error) {
	halfOrder, ok := curveHalfOrders[k.Curve]
	if !ok {
		return false, fmt.Errorf("curve not recognized [%s]", k.Curve)
	}
	return s.Cmp(halfOrder) != 1, nil
}

func toLowS(k *ecdsa.PublicKey, s *big.Int) (*big.Int, error) {
	lowS, err := isLowS(k, s)
	if err != nil {
		return nil, err
	}
	if !lowS {
		// Set s to N - s that will be then in the lower part of signature space
		// less or equal to half order
		s.Sub(k.Params().N, s)
		return s, nil
	}
	return s, nil
}
