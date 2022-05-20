package models

import (
	"fmt"
	"strings"
)

// Vote_<Vote.ElectionName>_<Vote.TxID>
const VOTE_KEY_TEMPLATE = "vote_%s_%s"

func NewVote(election *Election, candidate, txID string, nominations map[string]string) (*Vote, error) {
	vote := &Vote{
		ElectionName: election.Name,
		Candidate:    candidate,
		Nominations:  nominations,
		TxID:         txID,
	}

	if err := vote.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %s", err)
	}

	if err := vote.postValidate(election); err != nil {
		return nil, fmt.Errorf("validation error: %s", err)
	}

	return vote, nil
}

type Vote struct {
	ElectionName string `json:"electionName"`
	Candidate    string `json:"candidate"`
	// nomination => nominated
	Nominations map[string]string `json:"nominations"`
	TxID        string            `json:"txID"`
}

func (v *Vote) UniqueKey() string {
	return fmt.Sprintf(VOTE_KEY_TEMPLATE, v.ElectionName, v.TxID)
}

func (v *Vote) Validate() error {
	emptyFieldsErrMsgTemplate := "current fields are empty: [%s]"

	emptyFields := []string{}

	if v.ElectionName == "" {
		emptyFields = append(emptyFields, "electionName")
	}

	if v.Candidate == "" {
		emptyFields = append(emptyFields, "candidate")
	}

	if len(emptyFields) != 0 {
		return fmt.Errorf(emptyFieldsErrMsgTemplate, strings.Join(emptyFields, ", "))
	}

	return nil
}

// postValidate проверяет существует ли кандидат в текущем голосовании,
// а так же удаляет все номинации, которые не входят в голосовании
func (v *Vote) postValidate(election *Election) error {
	if _, ok := election.Candidates[v.Candidate]; !ok {
		return fmt.Errorf("vote candidate is not included in election")
	}

	for k := range v.Nominations {
		if _, ok := election.Nominations[k]; !ok {
			delete(v.Nominations, k)
		}
	}

	return nil
}
