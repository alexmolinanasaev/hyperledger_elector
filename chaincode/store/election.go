package store

import (
	"elector/chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type ElectionStore struct {
	store *store
}

func GetElectionStore(stub shim.ChaincodeStubInterface) *ElectionStore {
	return &ElectionStore{
		store: getStore(stub),
	}
}

func (s *ElectionStore) PutOne(election *models.Election) error {
	return s.store.putOne(election)
}

func (s *ElectionStore) PutMany(election []*models.Election) error {
	for _, d := range election {
		err := s.PutOne(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ElectionStore) GetOneByKey(key string) (*models.Election, error) {
	electionRaw, err := s.store.getOneByKey(key)
	if err != nil {
		return nil, fmt.Errorf("cannot get election by key: %s", err)
	}

	result := &models.Election{}
	if err := json.Unmarshal(electionRaw, result); err != nil {
		return nil, fmt.Errorf("cannot unmarshal election: %s", err)
	}

	return result, nil
}
