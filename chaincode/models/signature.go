package models

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"strings"

	"elector/chaincode/utils"
)

// signature_<Signature.SignedMessage(sha256 hex)>
const SIGNATURE_KEY_TEMPLATE = "signature_%s"

type Signature struct {
	ElectionName  string `json:"electionName,omitempty"`
	ElectorMSP    string `json:"electorMSP,omitempty"`
	SignedMessage string `json:"signedMessage"`
	MessageHash   []byte `json:"messageHash"`
	SignerPubKey  *ecdsa.PublicKey
}

func (s *Signature) UniqueKey() string {
	signatureHashBytes := sha256.Sum256([]byte(s.SignedMessage))

	return fmt.Sprintf(SIGNATURE_KEY_TEMPLATE, fmt.Sprintf("%x", signatureHashBytes))
}

func (s *Signature) Validate() error {
	// если хэш уже расчитан - не надо выполнять дальнейшую валидацию
	if len(s.MessageHash) == 32 {
		return nil
	}

	errMsgTemplate := "current fields are empty: [%s]"

	emptyFields := []string{}

	if s.ElectionName == "" {
		emptyFields = append(emptyFields, "electionName")
	}

	if s.ElectorMSP == "" {
		emptyFields = append(emptyFields, "electorMSP")
	}

	if s.SignedMessage == "" {
		emptyFields = append(emptyFields, "signedMessage")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(errMsgTemplate, strings.Join(emptyFields, ", "))
	}

	if s.SignerPubKey == nil {
		return fmt.Errorf("[INTERNAL] signer pub key not provided")
	}

	if ok := utils.VerifySignature(s.SignerPubKey, s.HashElectorPayload(), []byte(s.SignedMessage)); !ok {
		return fmt.Errorf("wrong signature")
	}

	return nil
}

func (s *Signature) HashElectorPayload() []byte {
	// если хэш уже расчитан - просто возвращаем его
	if s.MessageHash != nil {
		return s.MessageHash
	}

	s.MessageHash = utils.HashElectorPayload(s.ElectionName, s.ElectorMSP)

	return s.MessageHash
}
