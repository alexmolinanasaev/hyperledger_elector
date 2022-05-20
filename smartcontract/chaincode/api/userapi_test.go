package api_test

import (
	"elector/chaincode"
	"elector/chaincode/models"

	. "github.com/onsi/ginkgo"
	expectcc "github.com/s7techlab/cckit/testing/expect"
	// expectcc "github.com/s7techlab/cckit/testing/expect"
)

var _ = Describe("User API", func() {
	electionName := "Best Crypto Currency"

	vote1 := &models.Vote{
		ElectionName: electionName,
		Candidate:    "BTC",
		Nominations: map[string]string{
			"Most Stable":      "BTC",
			"Best Liquidity":   "BTC",
			"Best Perspective": "BTC",
			"44":               "MINA",
		},
	}

	// vote2 := &models.Vote{
	// 	ElectionName: electionName,
	// 	Candidate:    "BTC",
	// 	Nominations: map[string]string{
	// 		"Most Stable":      "Doggy",
	// 		"Best Liquidity":   "Doggy",
	// 		"Best Perspective": "nil",
	// 		"44":               "MINA",
	// 	},
	// }

	// wrongVoteTransient := voteToTransient(wrongVote, ELECTOR1_CORRECT_SIGNATURE)
	vote1Transient := voteToTransient(vote1, ELECTOR1_CORRECT_SIGNATURE)

	chaincode := &chaincode.SmartContract{}

	ctx, stub := prepMocksAsAdmin(vote1Transient)

	candidates := map[string]string{
		"BTC":   "Just HODL it",
		"USDT":  "Stable as democracy",
		"MINA":  "Should be 44!",
		"DOGGY": "Not scam",
	}

	nominations := map[string]string{
		"Most Stable":      "Minimal price jumps",
		"Best Liquidity":   "Biggest capital",
		"Best Perspective": "Coin you shoud HODL",
		"44":               "MINA 44!",
	}

	Context("Vote", func() {

		It("Fail because admin cannot vote", func() {

			stub.MockTransactionStart("vote")
			expectcc.ResponseError(chaincode.Vote(ctx), "admin cannot vote")
			stub.MockTransactionEnd("vote")
		})

		It("Fail because election does not exist", func() {
			wrongTransient := map[string][]byte{
				"electionName":  []byte("test"),
				"signedMessage": []byte(ELECTOR1_CORRECT_SIGNATURE),
			}

			ctx, stub := prepMocksAsElector1(wrongTransient)

			stub.MockTransactionStart("vote")
			expectcc.ResponseError(chaincode.Vote(ctx), "cannot vote to non existent election")
			stub.MockTransactionEnd("vote")
		})

		It("Success", func() {

			// создание выборов
			stub.MockTransactionStart("vote")
			expectcc.ResponseOk(chaincode.NewElection(ctx, electionName, candidates, nominations))
			// stub.MockTransactionEnd("vote")

			ctx, _ := prepMocksAsElector1(vote1Transient)

			// stub.MockTransactionStart("vote")
			expectcc.ResponseOk(chaincode.Vote(ctx))
			stub.MockTransactionEnd("vote")
		})

		// It("Fail because election does not exist", func() {
		// 	ctx, stub := prepMocksAsElector1(wrongVoteTransient)

		// 	stub.MockTransactionStart("vote")
		// 	expectcc.ResponseError(chaincode.Vote(ctx), "admin cannot vote")
		// 	stub.MockTransactionEnd("vote")
		// })
	})
})
