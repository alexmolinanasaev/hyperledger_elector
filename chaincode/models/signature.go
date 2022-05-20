package models

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	"elector/chaincode/utils"
)

// signature_<Signature.electionName>_<Signature.ElectorMSP(sha256 hex)>
const SIGNATURE_KEY_TEMPLATE = "signature_%s_%s"

func NewSignature(electionName, electorMSP, signedMessage, pubK string) (*Signature, error) {
	pubKey, err := utils.ExtractPubKeyFromCert([]byte(pubK))
	if err != nil {
		return nil, fmt.Errorf("[INTERNAL] cannot decode admin pub key")
	}

	signature := &Signature{
		ElectionName:  electionName,
		ElectorMSP:    electorMSP,
		SignedMessage: signedMessage,
		SignerPubKey:  pubKey,
	}

	if err := signature.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %s", err)
	}

	return signature, nil
}

type Signature struct {
	ElectionName  string `json:"electionName"`
	ElectorMSP    string `json:"electorMSP"`
	SignedMessage string `json:"signedMessage"`
	SignerPubKey  *ecdsa.PublicKey
}

func (s *Signature) UniqueKey() string {
	return fmt.Sprintf(SIGNATURE_KEY_TEMPLATE, s.ElectionName, fmt.Sprintf("%x", s.HashElectorPayload()))
}

func (s *Signature) Validate() error {
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
	return utils.HashElectorPayload(s.ElectionName, s.ElectorMSP)
}
