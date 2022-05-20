package api_test

import (
	"elector/chaincode/models"
	"encoding/json"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-samples/asset-transfer-private-data/chaincode-go/chaincode/mocks"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
)

const ADMIN_IDENTITY = "Org1MSP.eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT"
const ELECTOR1_MSP = "Org1MSP.eDUwOTo6Q049VXNlcjFAb3JnMS5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="
const ELECTOR2_MSP = "Org2MSP.eDUwOTo6Q049VXNlcjFAb3JnMi5leGFtcGxlLmNvbSxPVT1jbGllbnQsTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUzo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxPPW9yZzIuZXhhbXBsZS5jb20sTD1TYW4gRnJhbmNpc2NvLFNUPUNhbGlmb3JuaWEsQz1VUw=="

const ELECTOR1_CORRECT_SIGNATURE = "30440220677b73a526aea5f1356a0b4c339605cfab64ad70ccf76e6cdb1d47b8cc6f899d02202fdfcb4b52f4c5d39852e56323d95b26d347e5fd6fd8c95904c1c8ab6826b748"
const ELECTOR2_CORRECT_SIGNATURE = "304402206f6fc6438005b278a353c64fe7464be60d3088140a84a3b0516ae7c405cc3b3f022031acd15e477bc4d711aac2a8216bfd6696518c59bb02f2c0e1efe55f5a05f95d"

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API")
}

func prepMocksAsAdmin(transient map[string][]byte) (*mocks.TransactionContext, *shimtest.MockStub) {
	return prepMocks(ADMIN_IDENTITY, transient)
}

func prepMocksAsElector1(transient map[string][]byte) (*mocks.TransactionContext, *shimtest.MockStub) {
	return prepMocks(ELECTOR1_MSP, transient)
}

func prepMocksAsElector2(transient map[string][]byte) (*mocks.TransactionContext, *shimtest.MockStub) {
	return prepMocks(ELECTOR2_MSP, transient)
}

func prepMocks(identity string, transient map[string][]byte) (*mocks.TransactionContext, *shimtest.MockStub) {
	chaincodeStub := shimtest.NewMockStub("elector", cpaper_extended.NewCC())
	chaincodeStub.TransientMap = transient

	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	id := strings.Split(identity, ".")

	clientIdentity := &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns(id[0], nil)
	clientIdentity.GetIDReturns(id[1], nil)
	//set matching msp ID using peer shim env variable
	transactionContext.GetClientIdentityReturns(clientIdentity)
	return transactionContext, chaincodeStub
}

func voteToTransient(vote *models.Vote, signedMessage string) map[string][]byte {
	nominationsBytes, _ := json.Marshal(vote.Nominations)

	return map[string][]byte{
		"electionName":  []byte(vote.ElectionName),
		"candidate":     []byte(vote.Candidate),
		"signedMessage": []byte(signedMessage),
		"nominations":   nominationsBytes,
	}
}
