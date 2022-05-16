package models_test

import (
	"elector/chaincode/models"
	"elector/chaincode/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vote Model", func() {
	Context("Validation", func() {
		vote := &models.Vote{}

		It("Empty fields", func() {
			Expect(vote.Validate()).Should(MatchError("current fields are empty: [signature, candidate]"))
		})

		// так, как валидация сигнатуры тестируется в другом месте - тут главное проверить криптографическую часть сигнатуры
		// нужно проверить что вообще вызывается валидация сигнатуры
		It("Wrong signature validation error", func() {
			pubKey, err := utils.ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
			if err != nil {
				Fail("cannot extract pub key")
			}

			vote.Candidate = "BTC"
			vote.Signature = &models.Signature{
				ElectionName:  "Best Crypto Currency",
				ElectorMSP:    "Org2MSP",
				SignerPubKey:  pubKey,
				SignedMessage: WRONG_SIGNATURE,
			}

			Expect(vote.Validate()).Should(MatchError("signature validation error: wrong signature"))
		})

		It("Correct signature", func() {
			pubKey, err := utils.ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
			if err != nil {
				Fail("cannot extract pub key")
			}

			vote.Candidate = "BTC"
			vote.Signature = &models.Signature{
				ElectionName:  "Best Crypto Currency",
				ElectorMSP:    "Org2MSP",
				SignerPubKey:  pubKey,
				SignedMessage: CORRECT_SIGNATURE,
			}

			Expect(vote.Validate()).Should(MatchError("signature validation error: wrong signature"))
		})
	})
})
