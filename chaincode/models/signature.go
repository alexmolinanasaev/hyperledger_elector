package models

import (
	"crypto/sha256"
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
	MessageHash   string `json:"messageHash"`
}

func (s *Signature) UniqueKey() string {
	// очищаем прочие данные, чтобы они не хранились в блокчейне и голосование было анонимным
	signature := s.SignedMessage

	if s.MessageHash != "" {
		s.SignedMessage = ""
		s.ElectorMSP = ""
		s.ElectionName = ""
	}

	return fmt.Sprintf(SIGNATURE_KEY_TEMPLATE, signature)
}

func (s *Signature) Validate() error {
	// если хэш уже расчитан - не надо выполнять дальнейшую валидацию
	if s.MessageHash != "" {
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

	if ok := utils.CheckSignature(utils.GetAdminPub(), s.HashMessage(), []byte(s.SignedMessage)); !ok {
		return fmt.Errorf("wrong signature")
	}

	return nil
}

func (s *Signature) HashMessage() []byte {
	// ElectionName.ElectorMSP
	messageHash := "%s.%s"
	messageHash = fmt.Sprintf(messageHash, s.ElectionName, s.ElectorMSP)

	hashBytes := sha256.Sum256([]byte(messageHash))
	return hashBytes[:]
}
