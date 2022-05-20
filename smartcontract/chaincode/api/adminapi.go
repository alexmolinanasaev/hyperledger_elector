package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const ADMIN_IDENTITY = "Org1MSP.eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT"

type AdminAPI struct{}

func (api *AdminAPI) NewElection(ctx contractapi.TransactionContextInterface, name string, candidates, nominations map[string]string) peer.Response {
	if err := api.auth(ctx); err != nil {
		return shim.Error(err.Error())
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
	if err := api.auth(ctx); err != nil {
		return shim.Error(err.Error())
	}

	electionStore := store.GetElectionStore(ctx.GetStub())

	err := electionStore.CloseElectionByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (api *AdminAPI) auth(ctx contractapi.TransactionContextInterface) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("empty mspID: %s", err)
	}

	userID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("empty ID: %s", err)
	}

	if fmt.Sprintf("%s.%s", mspID, userID) != ADMIN_IDENTITY {
		return fmt.Errorf("can be called only by admin")
	}

	return nil
}
