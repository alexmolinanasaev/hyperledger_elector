package chaincode

import (
	"elector/chaincode/api"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.ContractInterface
	api.AdminAPI
}

type UserAPI struct{}
