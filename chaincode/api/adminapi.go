package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, name string, candidates, nominations map[string]string) error {
	election, err := models.NewElection(name, candidates, nominations)
	if err != nil {
		return err
	}

	electionStore := store.GetElectionStore(ctx.GetStub())

	return electionStore.PutOne(election)
}

func (api *AdminAPI) CloseElection(ctx contractapi.TransactionContextInterface, electionName string) error {
	electionStore := store.GetElectionStore(ctx.GetStub())

	return electionStore.CloseElectionByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
}
