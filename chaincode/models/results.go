package models

import "fmt"

type VotingResults struct {
	ElectionName string
	Candidates   map[string]int
	Nominations  map[string]map[string]int
}

func CountVotes(election *Election, votes []*Vote) (*VotingResults, error) {
	if !election.Closed {
		return nil, fmt.Errorf("election not closed yet")
	}

	candidates := make(map[string]int)
	nominations := make(map[string]map[string]int)

	for k := range election.Candidates {
		candidates[k] = 0
	}

	for k := range election.Nominations {
		if _, ok := election.Nominations[k]; !ok {
			continue
		}

		nominations[k] = make(map[string]int)

		for kk := range election.Candidates {
			nominations[k][kk] = 0
		}
	}

	for _, vote := range votes {
		candidates[vote.Candidate]++
		for k, v := range vote.Nominations {
			if _, ok := nominations[k]; !ok {
				continue
			}

			nominations[k][v]++
		}
	}

	return &VotingResults{
		ElectionName: election.Name,
		Candidates:   candidates,
		Nominations:  nominations,
	}, nil
}
