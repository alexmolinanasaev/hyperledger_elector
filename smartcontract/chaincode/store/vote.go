package store

import (
	"elector/chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type VoteStore struct {
	store *store
}

func GetVoteStore(stub shim.ChaincodeStubInterface) *VoteStore {
	return &VoteStore{
		store: getStore(stub),
	}
}

func (s *VoteStore) PutOne(vote *models.Vote) error {
	if err := vote.Validate(); err != nil {
		return fmt.Errorf("validation error: %s", err)
	}

	foundVote, err := s.GetOneByKey(vote.UniqueKey())
	if err != nil {
		return fmt.Errorf("cannot verify if already exist: %s", err)
	}

	if foundVote != nil {
		return fmt.Errorf("already exist")
	}

	return s.store.putOne(vote)
}

func (s *VoteStore) PutMany(vote []*models.Vote) error {
	for _, d := range vote {
		err := s.PutOne(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *VoteStore) GetOneByKey(key string) (*models.Vote, error) {
	voteRaw, err := s.store.getOneByKey(key)
	if err != nil {
		return nil, fmt.Errorf("cannot get vote by key: %s", err)
	}

	if voteRaw == nil {
		return nil, nil
	}

	result := &models.Vote{}
	if err := json.Unmarshal(voteRaw, result); err != nil {
		return nil, fmt.Errorf("cannot unmarshal vote: %s", err)
	}

	return result, nil
}

func (s *VoteStore) GetManyByElectionName(electionName string) ([]*models.Vote, error) {
	fromKey := fmt.Sprintf(models.VOTE_KEY_TEMPLATE, electionName, "")
	toKey := fmt.Sprintf(models.VOTE_KEY_TEMPLATE, electionName, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	votesRaw, err := s.store.getByKeyRange(fromKey, toKey)

	if err != nil {
		return nil, fmt.Errorf("cannot get votes: %s", err)
	}

	result := make([]*models.Vote, 0)

	for _, v := range votesRaw {
		vote := &models.Vote{}

		err := json.Unmarshal(v, vote)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal vote: %s", err)
		}

		result = append(result, vote)
	}

	return result, nil
}
