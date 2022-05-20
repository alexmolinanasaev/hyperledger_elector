package models_test

import (
	"elector/chaincode/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Results model", func() {
	expectedResults := &models.VotingResults{
		ElectionName: "Best Crypto Currency",
		Candidates: map[string]int{
			"BTC":   4,
			"USDT":  3,
			"MINA":  2,
			"DOGGY": 1,
		},
		Nominations: map[string]map[string]int{
			"Most Stable": {
				"BTC":   4,
				"USDT":  3,
				"MINA":  2,
				"DOGGY": 1,
			},
			"Best Liquidity": {
				"BTC":   1,
				"USDT":  2,
				"MINA":  3,
				"DOGGY": 4,
			},
			"Best Perspective": {
				"BTC":   5,
				"USDT":  0,
				"MINA":  0,
				"DOGGY": 5,
			},
			"44": {
				"BTC":   0,
				"USDT":  0,
				"MINA":  10,
				"DOGGY": 0,
			},
		},
	}

	votes := []*models.Vote{
		{
			Candidate: "BTC",
			Nominations: map[string]string{
				"Most Stable":      "BTC",
				"Best Liquidity":   "BTC",
				"Best Perspective": "BTC",
				"44":               "MINA",
			},
		},
		{
			Candidate: "BTC",
			Nominations: map[string]string{
				"Most Stable":      "BTC",
				"Best Liquidity":   "USDT",
				"Best Perspective": "BTC",
				"44":               "MINA",
			},
		},
		{
			Candidate: "BTC",
			Nominations: map[string]string{
				"Most Stable":      "BTC",
				"Best Liquidity":   "USDT",
				"Best Perspective": "BTC",
				"44":               "MINA",
			},
		},
		{
			Candidate: "BTC",
			Nominations: map[string]string{
				"Most Stable":      "BTC",
				"Best Liquidity":   "MINA",
				"Best Perspective": "BTC",
				"44":               "MINA",
			},
		},
		{
			Candidate: "USDT",
			Nominations: map[string]string{
				"Most Stable":      "USDT",
				"Best Liquidity":   "MINA",
				"Best Perspective": "BTC",
				"44":               "MINA",
			},
		},
		{
			Candidate: "USDT",
			Nominations: map[string]string{
				"Most Stable":      "USDT",
				"Best Liquidity":   "MINA",
				"Best Perspective": "DOGGY",
				"44":               "MINA",
			},
		},
		{
			Candidate: "USDT",
			Nominations: map[string]string{
				"Most Stable":      "USDT",
				"Best Liquidity":   "DOGGY",
				"Best Perspective": "DOGGY",
				"44":               "MINA",
			},
		},
		{
			Candidate: "MINA",
			Nominations: map[string]string{
				"Most Stable":      "MINA",
				"Best Liquidity":   "DOGGY",
				"Best Perspective": "DOGGY",
				"44":               "MINA",
			},
		},
		{
			Candidate: "MINA",
			Nominations: map[string]string{
				"Most Stable":      "MINA",
				"Best Liquidity":   "DOGGY",
				"Best Perspective": "DOGGY",
				"44":               "MINA",
			},
		},
		{
			Candidate: "DOGGY",
			Nominations: map[string]string{
				"Most Stable":      "DOGGY",
				"Best Liquidity":   "DOGGY",
				"Best Perspective": "DOGGY",
				"44":               "MINA",
			},
		},
	}

	election := &models.Election{
		Name: "Best Crypto Currency",
		Candidates: map[string]string{
			"BTC":   "Just HODL it",
			"USDT":  "Stable as democracy",
			"MINA":  "Should be 44!",
			"DOGGY": "Not scam",
		},
		Nominations: map[string]string{
			"Most Stable":      "Minimal price jumps",
			"Best Liquidity":   "Biggest capital",
			"Best Perspective": "Coin you shoud HODL",
			"44":               "MINA 44!",
		},
		Closed: true,
	}

	It("Fail because of opened election", func() {
		_, err := models.CountVotes(&models.Election{}, votes)
		Expect(err).Should(MatchError("election not closed yet"))
	})

	It("Sucess", func() {
		Expect(models.CountVotes(election, votes)).Should(Equal(expectedResults))
	})
})
