package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const ADMIN_IDENTITY = "Org1MSP"

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, name string, candidates, nominations map[string]string) peer.Response {
	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	userID, _ := ctx.GetClientIdentity().GetID()

	fmt.Println("mspID = ", mspID)
	fmt.Println("userID = ", userID)

	if fmt.Sprintf("%s%s", mspID, userID) != ADMIN_IDENTITY {
		return shim.Error("wrong identity")
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
	electionStore := store.GetElectionStore(ctx.GetStub())

	err := electionStore.CloseElectionByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
