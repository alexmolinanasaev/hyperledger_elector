package models

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"elector/chaincode/utils"
)

// signature_<Signature.SignedMessage>
const SIGNATURE_KEY_TEMPLATE = "vote_%s"

type Signature struct {
	ElectionName  string `json:"electionName,omitempty"`
	ElectorMSP    string `json:"electorMSP,omitempty"`
	SignedMessage string `json:"signedMessage,omitempty"`
	MessageHash   []byte `json:"messageHash"`
	SignerPubKey  *ecdsa.PublicKey
}

func (s *Signature) UniqueKey() string {
	// очищаем прочие данные, чтобы они не хранились в блокчейне и голосование было анонимным
	signature := s.SignedMessage

	if len(s.MessageHash) == 32 {
		s.SignedMessage = ""
		s.ElectorMSP = ""
		s.ElectionName = ""
	}

	return fmt.Sprintf(SIGNATURE_KEY_TEMPLATE, signature)
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

	if s.SignedMessage == "" {
		emptyFields = append(emptyFields, "message")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(errMsgTemplate, strings.Join(emptyFields, ", "))
	}

	if s.SignerPubKey == nil {
		return fmt.Errorf("[INTERNAL] signer pub key not provided")
	}

	if ok := utils.CheckSignature(s.SignerPubKey, s.HashElectorPayload(), []byte(s.SignedMessage)); !ok {
		return fmt.Errorf("wrong signature")
	}

	return nil
}

func (s *Signature) HashElectorPayload() []byte {
	// если хэш уже расчитан - просто возвращаем его
	if len(s.MessageHash) == 32 {
		return s.MessageHash
	}

	return utils.HashElectorPayload(s.ElectionName, s.ElectorMSP)
}
