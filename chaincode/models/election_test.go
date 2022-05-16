package models_test

import (
	"elector/chaincode/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Elector Model", func() {
	election := &models.Election{}

	Context("Validation", func() {
		It("validation failed", func() {
			Expect(election.Validate()).Should(MatchError("current fields are empty: [name, candidates]"))
		})

		It("Validation passed", func() {
			election.Name = "Best Crypto Currency"
			election.Candidates = []string{"BTC", "USDT", "MINA", "DOGGY"}
			election.Nominations = []string{"Most Stable", "Best Liquidity", "Best Perspective", "44"}

			Expect(election.Validate()).Should(Succeed())
		})
	})

	Context("Unique key", func() {
		It("Correct Unique Key", func() {
			Expect(election.UniqueKey()).Should(Equal("election_Best Crypto Currency"))
		})
	})

	Context("Close", func() {
		It("Is not closed yet", func() {
			Expect(election.Close()).Should(Succeed())
		})

		It("Is closed yet", func() {
			Expect(election.Close()).Should(MatchError("already closed"))
		})
	})
})
