package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"elector/chaincode/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type UserAPI struct{}

func (api *UserAPI) Vote(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	transient, err := stub.GetTransient()
	if err != nil {
		return fmt.Errorf("cannot get transient data: %s", err)
	}

	electionName, candidate, signedMessage, nominations := api.parseTransient(transient)

	electionStore := store.GetElectionStore(stub)

	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return fmt.Errorf("election not found: %s", err)
	}

	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	userID, _ := ctx.GetClientIdentity().GetID()

	electorMSP := fmt.Sprintf("%s%s", mspID, userID)

	signature, err := models.NewSignature(electionName, electorMSP, signedMessage, utils.ADMIN_PUB_KEY)
	if err != nil {
		return fmt.Errorf("signature error: %s", err)
	}

	signatureStore := store.GetSignatureStore(stub)
	s, err := signatureStore.GetOneByKey(signature.UniqueKey())
	if err != nil {
		return fmt.Errorf("cannot verify if signature is used: %s", err)
	}

	if s != nil {
		return fmt.Errorf("signature already used")
	}

	vote, err := models.NewVote(election, candidate, nominations)
	if err != nil {
		return err
	}

	voteStore := store.GetVoteStore(stub)
	if err := voteStore.PutOne(vote); err != nil {
		return fmt.Errorf("cannot save vote: %s", err)
	}

	if err := signatureStore.PutOne(signature); err != nil {
		return fmt.Errorf("cannot save signature: %s", err)
	}

	return nil
}

// func (api *UserAPI) GetResults(ctx contractapi.TransactionContextInterface, electionName string) {
// 	electionStore := store.GetElectionStore(stub)

// 	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
// 	if err != nil {
// 		return fmt.Errorf("election not found: %s", err)
// 	}
// }

func (api *UserAPI) GetElection(ctx contractapi.TransactionContextInterface, electionName string) (*models.Election, error) {
	return nil, nil
}

func (api *UserAPI) GetVotesCount(ctx contractapi.TransactionContextInterface, electionName string) {
}

func (api *UserAPI) parseTransient(transient map[string][]byte) (string, string, string, map[string]string) {
	// не будем проверять есть ли все поля в мапе. Валидатор сделает это за нас
	electionNameBytes := transient["electionName"]
	electionName := string(electionNameBytes)

	candidateBytes := transient["candidate"]
	candidate := string(candidateBytes)

	signedMessageBytes := transient["signedMessage"]
	signedMessage := fmt.Sprintf("%x", signedMessageBytes)

	nominationsBytes := transient["nominations"]
	nominationsRaw := make(map[string][]byte)
	json.Unmarshal(nominationsBytes, &nominationsRaw)
	nominations := make(map[string]string)
	for k, v := range nominations {
		nominations[k] = string(v)
	}

	return electionName, candidate, signedMessage, nominations
}