package models

import (
	"fmt"
	"strings"
)

// Vote_<Vote.Election.Name>
const VOTE_KEY_TEMPLATE = "vote_%s"

type Vote struct {
	Signature   *Signature        `json:"signature,omitempty"`
	Candidate   string            `json:"candidate"`
	Nominations map[string]string `json:"nominations"`
}

func (v *Vote) UniqueKey() string {
	return fmt.Sprintf(VOTE_KEY_TEMPLATE, v.Signature.ElectionName)
}

func (v *Vote) Validate() error {
	emptyFieldsErrMsgTemplate := "current fields are empty: [%s]"

	emptyFields := []string{}

	if v.Signature == nil {
		emptyFields = append(emptyFields, "signature")
	}

	if v.Candidate == "" {
		emptyFields = append(emptyFields, "candidate")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(emptyFieldsErrMsgTemplate, strings.Join(emptyFields, ", "))
	}

	if err := v.Signature.Validate(); err != nil {
		return fmt.Errorf("signature validation error: %s", err)
	}

	return nil
}
