package api

import (
	"elector/chaincode/models"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, election models.Election) error {
	return nil
}

func (api *AdminAPI) CloseElection(ctx contractapi.TransactionContextInterface, electionName string) error {
	return nil
}
