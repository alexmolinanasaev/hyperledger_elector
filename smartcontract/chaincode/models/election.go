package models

import (
	"fmt"
	"strings"
)

// election_<Election.Name>
const ELECTION_KEY_TEMPLATE = "election_%s"

func NewElection(name string, candidates, nominations map[string]string) (*Election, error) {
	election := &Election{
		Name:        name,
		Candidates:  candidates,
		Nominations: nominations,
		Closed:      false,
	}

	if err := election.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %s", err)
	}

	return election, nil
}

type Election struct {
	Name string `json:"name"`
	// candidate => description
	Candidates map[string]string `json:"candidates"`
	// nomination => description
	Nominations map[string]string `json:"nominations"`
	Closed      bool              `json:"closed"`
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

	// в голосовании могут отсутствовать номинации

	if len(emptyFields) != 0 {
		return fmt.Errorf(errMsgTemplate, strings.Join(emptyFields, ", "))
	}

	return nil
}

func (e *Election) Close() error {
	if e.Closed {
		return fmt.Errorf("already closed")
	}

	e.Closed = true

	return nil
}
