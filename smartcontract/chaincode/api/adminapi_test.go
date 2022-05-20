package api_test

import (
	"elector/chaincode"

	. "github.com/onsi/ginkgo"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

var _ = Describe("Admin API", func() {
	chaincode := &chaincode.SmartContract{}

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

	electionName := "Best Crypto Currency"

	ctx, stub := prepMocksAsAdmin(nil)
	Context("New Election", func() {
		Context("Admin role", func() {
			It("Sucess", func() {
				stub.MockTransactionStart("save election")
				expectcc.ResponseOk(chaincode.NewElection(ctx, electionName, candidates, nominations))
				stub.MockTransactionEnd("save election")
			})

			It("Fail because already exist", func() {
				stub.MockTransactionStart("save election")
				expectcc.ResponseError(chaincode.NewElection(ctx, electionName, candidates, nominations), "already exist")
				stub.MockTransactionEnd("save election")
			})
		})

		It("Fail because of wrong identity", func() {
			ctx, stub := prepMocksAsElector1(nil)
			stub.MockTransactionStart("save election")
			expectcc.ResponseError(chaincode.NewElection(ctx, electionName, candidates, nominations), "can be called only by admin")
			stub.MockTransactionEnd("save election")
		})
	})

	Context("Close Election", func() {
		Context("Admin role", func() {
			It("Sucess", func() {
				stub.MockTransactionStart("close election")
				expectcc.ResponseOk(chaincode.CloseElection(ctx, electionName))
				stub.MockTransactionEnd("close election")
			})

			It("Fail because already closed", func() {
				stub.MockTransactionStart("close election")
				expectcc.ResponseError(chaincode.CloseElection(ctx, electionName), "already closed")
				stub.MockTransactionEnd("close election")
			})

			It("Fail because not exist", func() {
				stub.MockTransactionStart("close election")
				expectcc.ResponseError(chaincode.CloseElection(ctx, "nil"), "non existent election cannot be closed")
				stub.MockTransactionEnd("close election")
			})

			It("Fail because of wrong identity", func() {
				ctx, stub := prepMocksAsElector1(nil)
				stub.MockTransactionStart("close election")
				expectcc.ResponseError(chaincode.CloseElection(ctx, electionName), "can be called only by admin")
				stub.MockTransactionEnd("close election")
			})
		})
	})
})
