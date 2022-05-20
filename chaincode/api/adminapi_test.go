package api_test

import (
	"elector/chaincode"

	. "github.com/onsi/ginkgo"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

var _ = Describe("contract API", func() {
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

	Context("Admin API", func() {
		ctx, stub := prepMocksAsAdmin()
		Context("New Election", func() {
			Context("Admin role", func() {
				It("Sucess", func() {
					stub.MockTransactionStart("save signature")
					expectcc.ResponseOk(chaincode.NewElection(ctx, "Best Crypto Currency", candidates, nominations))
					stub.MockTransactionEnd("save signature")
				})

				It("Fail because already exist", func() {
					stub.MockTransactionStart("save signature")
					expectcc.ResponseError(chaincode.NewElection(ctx, "Best Crypto Currency", candidates, nominations), "already exist")
					stub.MockTransactionEnd("save signature")
				})
			})

			It("Fail because of wrong identity", func() {
				ctx, stub := prepMocksAsElector1()
				stub.MockTransactionStart("save signature")
				expectcc.ResponseError(chaincode.NewElection(ctx, "Best Crypto Currency", candidates, nominations), "can be called only by admin")
				stub.MockTransactionEnd("save signature")
			})
		})

		Context("Close Election", func() {
			Context("Admin role", func() {
				It("Sucess", func() {
					stub.MockTransactionStart("save signature")
					expectcc.ResponseOk(chaincode.CloseElection(ctx, "Best Crypto Currency"))
					stub.MockTransactionEnd("save signature")
				})

				It("Fail because already closed", func() {
					stub.MockTransactionStart("save signature")
					expectcc.ResponseError(chaincode.CloseElection(ctx, "Best Crypto Currency"), "already closed")
					stub.MockTransactionEnd("save signature")
				})

				It("Fail because not exist", func() {
					stub.MockTransactionStart("save signature")
					expectcc.ResponseError(chaincode.CloseElection(ctx, "bad election name"), "non existent election cannot be closed")
					stub.MockTransactionEnd("save signature")
				})
			})

			It("Fail because of wrong identity", func() {
				ctx, stub := prepMocksAsElector1()
				stub.MockTransactionStart("save signature")
				expectcc.ResponseError(chaincode.CloseElection(ctx, "Best Crypto Currency"), "can be called only by admin")
				stub.MockTransactionEnd("save signature")
			})
		})
	})
})
