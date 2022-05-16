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
