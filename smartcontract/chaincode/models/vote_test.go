package models_test

import (
	"elector/chaincode/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vote Model", func() {
	Context("Validation", func() {
		vote := &models.Vote{}

		It("Empty fields", func() {
			Expect(vote.Validate()).Should(MatchError("current fields are empty: [electionName, candidate]"))
		})

		It("Sucess", func() {
			vote.ElectionName = "Best Crypto Currency"
			vote.Candidate = "BTC"
			Expect(vote.Validate()).Should(Succeed())
		})
	})
})
