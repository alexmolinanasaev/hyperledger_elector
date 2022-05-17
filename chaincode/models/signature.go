package models

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
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
	if s.MessageHash == nil {
		signedMessageBytes, err := hex.DecodeString(s.SignedMessage)
		if err != nil {
			return "" // чтобы удовлетворять интерфейсу store.storeable м ыпросто вернем пустой ключ
		}

		signatureHashBytes := sha256.Sum256(signedMessageBytes)
		s.MessageHash = signatureHashBytes[:]
	}

	return fmt.Sprintf(SIGNATURE_KEY_TEMPLATE, fmt.Sprintf("%x", s.MessageHash))
}

func (s *Signature) Validate() error {
	// если хэш уже расчитан - не надо выполнять дальнейшую валидацию
	if s.MessageHash != nil {
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

	signedMessageBytes, err := hex.DecodeString(s.SignedMessage)
	if err != nil {
		return fmt.Errorf("cannot parse signature hex: %s", err)
	}

	if ok := utils.VerifySignature(s.SignerPubKey, s.HashElectorPayload(), signedMessageBytes); !ok {
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
