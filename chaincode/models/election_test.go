package models

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestElectorModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Elector Model Suite")
}

var _ = Describe("Elector Model", func() {

	Context("Validation", func() {
		election := &Election{}
		It("validation failed", func() {
			Expect(election.Validate()).Should(MatchError("current fields are empty: [name, candidates]"))
		})

		It("Validation passed", func() {
			election.Name = "Best Crypto Currency"
			election.Candidates = []string{"BTC", "USDT", "MINA", "DOGGY"}
			election.Nominations = []string{"Most Stable", "Best Liquidity", "Best Perspective"}

			Expect(election.Validate()).Should(BeNil())
		})
	})

	Context("Validation passed", func() {
	})

	Context("Unique Key", func() {
		election := &Election{
			Name: "Best Crypto Currency",
		}

		It("Correct Unique Key", func() {
			Expect(election.UniqueKey()).Should(Equal("election_Best Crypto Currency"))
		})
	})

	Context("Close", func() {
		election := &Election{}

		It("Is not closed yet", func() {
			Expect(election.Close()).Should(BeNil())
		})

		It("Is closed yet", func() {
			Expect(election.Close()).Should(MatchError("already closed"))
		})
	})
})
