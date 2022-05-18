package main

import "elector/chaincode/utils"

// import (
// 	"elector/chaincode"
// 	"log"

// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// )

func main() {
	// salairSmartContract, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	// if err != nil {
	// 	log.Panicf("error creating chaincode: %v", err)
	// }
	// if err := salairSmartContract.Start(); err != nil {
	// 	log.Panicf("error starting chaincode: %v", err)
	// }

	utils.SignTestMessage()
}
