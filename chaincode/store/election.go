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
	if err := election.Validate(); err != nil {
		return fmt.Errorf("validation error: %s", err)
	}

	// нельзя вставить голосование если такое уже есть
	foundElection, err := s.GetOneByKey(election.UniqueKey())
	if err != nil {
		return fmt.Errorf("cannot verify if already exist: %s", err)
	}

	if foundElection != nil {
		return fmt.Errorf("already exist")
	}

	// нельзя сохранить закрытое голосование
	election.Closed = false
	return s.store.putOne(election)
}

func (s *ElectionStore) CloseElectionByKey(key string) error {
	// сущность голосования нельзя изменить после создания. Голосование можно только закрыть
	// если закрываемое голосование не существует, то и закрыть мы не можем
	foundElection, err := s.GetOneByKey(key)
	if err != nil || foundElection == nil {
		return fmt.Errorf("non existent election cannot be closed")
	}

	if err := foundElection.Close(); err != nil {
		return err
	}

	return s.store.putOne(foundElection)
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

	if electionRaw == nil {
		return nil, nil
	}

	result := &models.Election{}
	if err := json.Unmarshal(electionRaw, result); err != nil {
		return nil, fmt.Errorf("cannot unmarshal election: %s", err)
	}

	return result, nil
}
