package store

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type storeable interface {
	UniqueKey() string
}

type store struct {
	stub shim.ChaincodeStubInterface
}

func getStore(stub shim.ChaincodeStubInterface) *store {
	return &store{
		stub: stub,
	}
}

func (s *store) putOne(data storeable) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal data: %s", err)
	}

	err = s.stub.PutState(data.UniqueKey(), jsonData)
	if err != nil {
		return err
	}
	return nil
}

// func (s *store) putMany(data []storeable) error {
// 	for _, d := range data {
// 		err := s.putOne(d)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (s *store) getOneByKey(key string) ([]byte, error) {
	return s.stub.GetState(key)
}

func (s *store) getByKeyRange(startKey, endKey string) ([][]byte, error) {
	resultsIterator, err := s.stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := make([][]byte, 0)

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		results = append(results, queryResponse.Value)
	}

	return results, nil
}
