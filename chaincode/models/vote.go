package models

import (
	"fmt"
	"strings"
)

// Vote_<Vote.Election.Name>
const VOTE_KEY_TEMPLATE = "vote_%s"

type Vote struct {
	ElectionName string `json:"electionName"`
	Candidate    string `json:"candidate"`
	// nomination => nominated
	Nominations map[string]string `json:"nominations"`
}

func (v *Vote) UniqueKey() string {
	return fmt.Sprintf(VOTE_KEY_TEMPLATE, v.ElectionName)
}

func (v *Vote) Validate() error {
	emptyFieldsErrMsgTemplate := "current fields are empty: [%s]"

	emptyFields := []string{}

	if v.Candidate == "" {
		emptyFields = append(emptyFields, "candidate")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(emptyFieldsErrMsgTemplate, strings.Join(emptyFields, ", "))
	}

	return nil
}
