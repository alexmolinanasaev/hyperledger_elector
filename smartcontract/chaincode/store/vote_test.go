package store_test

import (
	"elector/chaincode/models"
	"elector/chaincode/store"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
)

var _ = Describe("Vote store", func() {
	electorChaincode := shimtest.NewMockStub(`elector`, cpaper_extended.NewCC())

	voteStore := store.GetVoteStore(electorChaincode)

	vote := &models.Vote{}
	Context("Put one", func() {
		It("Empty fields validation error", func() {
			Expect(voteStore.PutOne(vote)).Should(MatchError("validation error: current fields are empty: [electionName, candidate]"))
		})

		It("Success", func() {
			vote.ElectionName = "Best Crypto Currency"
			vote.Candidate = "BTC"

			electorChaincode.MockTransactionStart("save vote")
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			electorChaincode.MockTransactionEnd("save vote")
		})

		It("Already exist", func() {
			electorChaincode.MockTransactionStart("save vote")
			Expect(voteStore.PutOne(vote)).Should(MatchError("already exist"))
			electorChaincode.MockTransactionEnd("save vote")
		})
	})

	Context("Get one", func() {
		It("Not found", func() {
			electorChaincode.MockTransactionStart("get vote")
			v, err := voteStore.GetOneByKey("wrongKey")
			Expect(v).Should(BeNil())
			Expect(err).Should(BeNil())
			electorChaincode.MockTransactionEnd("get vote")
		})

		It("Success", func() {
			electorChaincode.MockTransactionStart("get vote")
			Expect(voteStore.GetOneByKey(vote.UniqueKey())).Should(Equal(vote))
			electorChaincode.MockTransactionEnd("get vote")
		})
	})

	Context("Get may", func() {
		It("Success", func() {
			electorChaincode.MockTransactionStart("save vote")
			vote.TxID = "1"
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			vote.TxID = "2"
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			vote.TxID = "3"
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			vote.TxID = "4"
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			vote.TxID = "5"
			Expect(voteStore.PutOne(vote)).Should(Succeed())
			vote.TxID = "6"
			Expect(voteStore.PutOne(vote)).Should(Succeed())

			votes, err := voteStore.GetManyByElectionName("Best Crypto Currency")

			Expect(err).Should(BeNil())
			Expect(len(votes)).Should(Equal(7))

			electorChaincode.MockTransactionEnd("save vote")
		})
	})
})
