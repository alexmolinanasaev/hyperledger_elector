package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const ADMIN_PUB_KEY string = `-----BEGIN CERTIFICATE-----
MIICKTCCAc+gAwIBAgIQTUFZM0uHkjSVzbXSGfdXAjAKBggqhkjOPQQDAjBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMjA1MTUxMTI2MDBaFw0zMjA1MTIxMTI2MDBa
MGsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMQ4wDAYDVQQLEwVhZG1pbjEfMB0GA1UEAwwWQWRtaW5Ab3Jn
MS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABCVcwwpmCcxU
oAdwCBJPr3kBaNPGpqFCzYXZ/zv0RNGeBm0Z07bA07lwhNZ6HtWIHjhRXbkKYM4i
49ctbnCxKS+jTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1Ud
IwQkMCKAIJfLMIp/JAvDlsfhsnFqpClhPpt0IVJwFlZkSnh13wxcMAoGCCqGSM49
BAMCA0gAMEUCIQC0cwSVvkx8oTh/87dERe7lnYDl5ZPyBuyBA5dSWs7s/AIgR889
qwRxuMGZG6KsLXw4P9zdFccUKEIIweVuOMkO1J0=
-----END CERTIFICATE-----`

func VerifySignature(pub *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	return ecdsa.VerifyASN1(pub, hash, sig)
}

func HashElectorPayload(electionName, electorMSP string) []byte {
	// ElectionName.ElectorMSP
	messageHash := "%s.%s"
	messageHash = fmt.Sprintf(messageHash, electionName, electorMSP)

	hashBytes := sha256.Sum256([]byte(messageHash))

	return hashBytes[:]
}

func ExtractPubKeyFromCert(certPEM []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*ecdsa.PublicKey), nil
}
