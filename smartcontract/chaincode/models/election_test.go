package models_test

import (
	"elector/chaincode/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Election Model", func() {
	election := &models.Election{}

	Context("Validation", func() {
		It("Empty fields", func() {
			Expect(election.Validate()).Should(MatchError("current fields are empty: [name, candidates]"))
		})

		It("Success", func() {
			election.Name = "Best Crypto Currency"
			election.Candidates = map[string]string{
				"BTC":   "Just HODL it",
				"USDT":  "Stable as democracy",
				"MINA":  "Should be 44!",
				"DOGGY": "Not scam",
			}
			election.Nominations = map[string]string{
				"Most Stable":      "Minimal price jumps",
				"Best Liquidity":   "Biggest capital",
				"Best Perspective": "Coin you shoud HODL",
				"44":               "MINA 44!",
			}

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
