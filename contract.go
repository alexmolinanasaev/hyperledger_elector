package main

import (
	"elector/chaincode/utils"
	"fmt"
	"log"
)

// import (
// 	"log"

// 	"elector/chaincode"

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

	signer, err := utils.NewSigner([]byte(utils.ADMIN_PRIV_KEY))
	if err != nil {
		log.Panic(err)
	}

	signature, err := signer.SignElectorPayload("Best Crypto Currency", "Org2MSP")
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println(len(signature))
	// fmt.Println(len(fmt.Sprintf("%x", signature)))

	// fmt.Printf("%x", signature)

	pub, err := utils.ExtractPubKeyFromCert([]byte(utils.ADMIN_PUB_KEY))
	if err != nil {
		log.Panic(err)
	}

	payload := utils.HashElectorPayload("Best Crypto Currency", "Org2MSP")

	ok := utils.VerifySignature(pub, payload, signature)
	fmt.Println(ok)
}
