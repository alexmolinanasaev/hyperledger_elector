package api

import (
	"elector/chaincode/models"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, name string, candidates, nominations map[string]string) error {
	election, err := models.NewElection(name, candidates, nominations)
	return nil
}

func (api *AdminAPI) CloseElection(ctx contractapi.TransactionContextInterface, electionName string) error {
	return nil
}
