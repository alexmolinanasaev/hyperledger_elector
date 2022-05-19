package main

import (
	"elector/chaincode"
	"elector/chaincode/api"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	chaincode := &chaincode.SmartContract{
		AdminAPI: api.AdminAPI{},
		UserAPI:  api.UserAPI{},
	}

	salairSmartContract, err := contractapi.NewChaincode(chaincode)
	if err != nil {
		log.Panicf("error creating chaincode: %v", err)
	}

	if err := salairSmartContract.Start(); err != nil {
		log.Panicf("error starting chaincode: %v", err)
	}
}
