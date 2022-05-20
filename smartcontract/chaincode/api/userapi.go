package api

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"elector/chaincode/utils"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type UserAPI struct{}

func (api *UserAPI) Vote(ctx contractapi.TransactionContextInterface) peer.Response {
	isAdmin, err := utils.IsAdmin(ctx)
	if err != nil {
		return shim.Error(fmt.Sprintf("identification error: %s", err))
	}

	if isAdmin {
		return shim.Error("admin cannot vote")
	}

	stub := ctx.GetStub()

	transient, err := stub.GetTransient()
	if err != nil {
		return shim.Error(fmt.Sprintf("cannot get transient data: %s", err))
	}

	electionName, candidate, signedMessage, nominations := api.parseTransient(transient)

	electionStore := store.GetElectionStore(stub)

	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(fmt.Sprintf("election not found: %s", err))
	}

	if election == nil {
		return shim.Error("cannot vote to non existent election")
	}

	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	userID, _ := ctx.GetClientIdentity().GetID()

	electorMSP := fmt.Sprintf("%s.%s", mspID, userID)

	signature, err := models.NewSignature(electionName, electorMSP, signedMessage, utils.ADMIN_PUB_KEY)
	if err != nil {
		return shim.Error(fmt.Sprintf("signature error: %s", err))
	}

	signatureStore := store.GetSignatureStore(stub)
	s, err := signatureStore.GetOneByKey(signature.UniqueKey())
	if err != nil {
		return shim.Error(fmt.Sprintf("cannot verify if signature is used: %s", err))
	}

	if s != nil {
		return shim.Error("signature already used")
	}

	vote, err := models.NewVote(election, candidate, stub.GetTxID(), nominations)
	if err != nil {
		return shim.Error(err.Error())
	}

	voteStore := store.GetVoteStore(stub)
	if err := voteStore.PutOne(vote); err != nil {
		return shim.Error(fmt.Sprintf("cannot save vote: %s", err))
	}

	if err := signatureStore.PutOne(signature); err != nil {
		return shim.Error(fmt.Sprintf("cannot save signature: %s", err))
	}

	return shim.Success(nil)
}

func (api *UserAPI) GetResults(ctx contractapi.TransactionContextInterface, electionName string) peer.Response {
	stub := ctx.GetStub()

	electionStore := store.GetElectionStore(stub)

	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(fmt.Sprintf("election not found: %s", err))
	}

	voteStore := store.GetVoteStore(stub)

	votes, err := voteStore.GetManyByElectionName(electionName)
	if err != nil {
		return shim.Error(err.Error())
	}

	votingResults, err := models.CountVotes(election, votes)
	if err != nil {
		return shim.Error(err.Error())
	}

	votingResultsBytes, err := json.Marshal(votingResults)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(votingResultsBytes)
}

func (api *UserAPI) GetElection(ctx contractapi.TransactionContextInterface, electionName string) peer.Response {
	electionStore := store.GetElectionStore(ctx.GetStub())

	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(err.Error())
	}

	electionBytes, err := json.Marshal(election)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(electionBytes)
}

func (api *UserAPI) GetVotes(ctx contractapi.TransactionContextInterface, electionName string) peer.Response {
	electionStore := store.GetElectionStore(ctx.GetStub())

	election, err := electionStore.GetOneByKey(fmt.Sprintf(models.ELECTION_KEY_TEMPLATE, electionName))
	if err != nil {
		return shim.Error(err.Error())
	}

	if !election.Closed {
		return shim.Error("cannot reveal election votes: election is not closed")
	}

	voteStore := store.GetVoteStore(ctx.GetStub())

	votes, err := voteStore.GetManyByElectionName(electionName)
	if err != nil {
		return shim.Error(err.Error())
	}

	response, err := json.Marshal(votes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(response)
}

func (api *UserAPI) GetVotesCount(ctx contractapi.TransactionContextInterface, electionName string) peer.Response {
	voteStore := store.GetVoteStore(ctx.GetStub())

	votes, err := voteStore.GetManyByElectionName(electionName)
	if err != nil {
		return shim.Error(err.Error())
	}

	response, err := json.Marshal(len(votes))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(response)
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
