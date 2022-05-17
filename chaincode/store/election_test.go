package store_test

import (
	"elector/chaincode/models"
	"elector/chaincode/store"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
)

var _ = Describe("Election store", func() {
	electorChaincode := shimtest.NewMockStub(`elector`, cpaper_extended.NewCC())

	electionStore := store.GetElectionStore(electorChaincode)

	election := &models.Election{}

	Context("Put one", func() {
		// главное проверить что валидация не прошла - неважно по какой причине
		It("Failed validation", func() {
			Expect(electionStore.PutOne(election)).Should(MatchError("validation error: current fields are empty: [name, candidates]"))
		})

		It("Success", func() {
			election = &models.Election{
				Name:        "Best Crypto Currency",
				Candidates:  []string{"BTC", "USDT", "MINA", "DOGGY"},
				Nominations: []string{"Most Stable", "Best Liquidity", "Best Perspective", "44"},
			}

			electorChaincode.MockTransactionStart("save election")
			Expect(electionStore.PutOne(election)).Should(Succeed())
			electorChaincode.MockTransactionEnd("save election")
		})

		It("Already exist", func() {
			electorChaincode.MockTransactionStart("save election")
			Expect(electionStore.PutOne(election)).Should(MatchError("already exist"))
			electorChaincode.MockTransactionEnd("save election")
		})
	})

	Context("Get one", func() {
		It("Not found", func() {
			electorChaincode.MockTransactionStart("get election")

			e, err := electionStore.GetOneByKey("wrongKey")
			Expect(e).Should(BeNil())
			Expect(err).Should(BeNil())

			electorChaincode.MockTransactionEnd("get election")
		})

		It("Success", func() {
			electorChaincode.MockTransactionStart("get election")
			Expect(electionStore.GetOneByKey(election.UniqueKey())).Should(Equal(election))
			electorChaincode.MockTransactionEnd("get election")
		})
	})

	Context("Close", func() {
		It("Success", func() {
			electorChaincode.MockTransactionStart("close election")
			Expect(electionStore.CloseElection(election)).Should(Succeed())
			electorChaincode.MockTransactionEnd("close election")
		})

		It("Already closed", func() {
			electorChaincode.MockTransactionStart("close election")
			Expect(electionStore.CloseElection(election)).Should(MatchError("already closed"))
			electorChaincode.MockTransactionEnd("close election")
		})

		It("Closing non existent election", func() {
			e := &models.Election{
				Name:       "non existent",
				Candidates: []string{"1"},
			}

			electorChaincode.MockTransactionStart("close election")
			Expect(electionStore.CloseElection(e)).Should(MatchError("non existent election cannot be closed"))
			electorChaincode.MockTransactionEnd("close election")
		})

		It("Closing election does not mumate other fields", func() {
			Expect(electionStore.GetOneByKey(election.UniqueKey())).Should(Equal(election))
		})
	})
})
