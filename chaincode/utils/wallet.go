package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
)

// TODO: пока только для теста
func SignTestMessage() {
	signer, err := NewSigner([]byte(ADMIN_PRIV_KEY))
	if err != nil {
		log.Panic(err)
	}

	signature, err := signer.SignElectorPayload("Best Crypto Currency", "Org2MSP")
	if err != nil {
		log.Panic(err)
	}

	pub, err := ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
	if err != nil {
		log.Panic(err)
	}

	payload := HashElectorPayload("Best Crypto Currency", "Org2MSP")

	ok := VerifySignature(pub, payload, signature)
	fmt.Println(ok)

	signatureHex := fmt.Sprintf("%x", signature)
	signatureHashBytes := sha256.Sum256([]byte(signatureHex))
	signatureHashHex := fmt.Sprintf("%x", signatureHashBytes[:])

	fmt.Println(signatureHashHex)
	fmt.Println(signatureHex)
}
