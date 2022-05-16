package store

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type storeable interface {
	UniqueKey() string
	Validate() error
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
	err := data.Validate()
	if err != nil {
		return fmt.Errorf("validation error: %s", err)
	}

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

func (s *store) putMany(data []storeable) error {
	for _, d := range data {
		err := s.putOne(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) getOneByKey(key string) ([]byte, error) {
	return s.stub.GetState(key)
}
