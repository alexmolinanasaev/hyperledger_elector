package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"elector/chaincode/utils"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, name string, candidates, nominations map[string]string) peer.Response {
	isAdmin, err := utils.IsAdmin(ctx)
	if err != nil {
		return shim.Error(fmt.Sprintf("identification error: %s", err))
	}

	if !isAdmin {
		return shim.Error("can be called only by admin")
	}

	election, err := models.NewElection(name, candidates, nominations)
	if err != nil {
		return shim.Error(err.Error())
	}

	electionStore := store.GetElectionStore(ctx.GetStub())

	err = electionStore.PutOne(election)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (api *AdminAPI) CloseElection(ctx contractapi.TransactionContextInterface, electionName string) peer.Response {
	isAdmin, err := utils.IsAdmin(ctx)
	if err != nil {
		return shim.Error(fmt.Sprintf("identification error: %s", err))
	}

	if !isAdmin {
		return shim.Error("can be called only by admin")
	}

	electionStore := store.GetElectionStore(ctx.GetStub())

	err = electionStore.CloseElectionByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
