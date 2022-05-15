package models

import (
	"fmt"
	"strings"
)

// election_<Election.Name>
const ELECTION_KEY_TEMPLATE = "election_%s"

type Election struct {
	Name        string   `json:"name"`
	Candidates  []string `json:"candidates"`
	Nominations []string `json:"nominations"`
	Closed      bool     `json:"closed"`
}

func (e *Election) UniqueKey() string {
	return fmt.Sprintf(ELECTION_KEY_TEMPLATE, e.Name)
}

func (e *Election) Validate() error {
	errMsgTemplate := "current fields are empty: [%s]"

	emptyFields := []string{}

	if e.Name == "" {
		emptyFields = append(emptyFields, "name")
	}

	if len(e.Candidates) == 0 {
		emptyFields = append(emptyFields, "candidates")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(errMsgTemplate, strings.Join(emptyFields, ", "))
	}

	// в голосовании могут отсутствовать номинации

	return nil
}

func (e *Election) Close() error {
	if e.Closed {
		return fmt.Errorf("already closed")
	}

	e.Closed = true

	return nil
}
