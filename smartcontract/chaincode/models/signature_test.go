package models_test

import (
	"elector/chaincode/models"
	"elector/chaincode/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signature Model", func() {
	Context("Validation", func() {
		signature := &models.Signature{}

		It("Empty fields", func() {
			Expect(signature.Validate()).Should(MatchError("current fields are empty: [electionName, electorMSP, signedMessage]"))
		})

		It("No pub key", func() {
			signature.ElectionName = "Best Crypto Currency"
			signature.ElectorMSP = "Org2MSP"
			signature.SignedMessage = WRONG_SIGNATURE

			Expect(signature.Validate()).Should(MatchError("[INTERNAL] signer pub key not provided"))
		})

		It("Wrong signature", func() {
			pubKey, err := utils.ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
			if err != nil {
				Fail("cannot extract pub key")
			}

			signature.SignerPubKey = pubKey

			Expect(signature.Validate()).Should(MatchError("wrong signature"))
		})

		It("Wrong election name", func() {
			signature.ElectionName = "Worst Crypto Currency"
			signature.ElectorMSP = "Org2MSP"
			signature.SignedMessage = CORRECT_SIGNATURE

			Expect(signature.Validate()).Should(MatchError("wrong signature"))
		})

		It("Wrong elector MSP", func() {
			signature.ElectionName = "Best Crypto Currency"
			signature.ElectorMSP = "Org3MSP"
			signature.SignedMessage = CORRECT_SIGNATURE

			Expect(signature.Validate()).Should(MatchError("wrong signature"))
		})

		It("Correct signature", func() {
			s, _ := models.NewSignature("Best Crypto Currency", "Org2MSP", CORRECT_SIGNATURE, ADMIN_PUB_KEY)
			Expect(s.Validate()).Should(Succeed())
		})
	})

	Context("Hash elector payload", func() {
		signature := &models.Signature{
			ElectionName: "Best Crypto Currency",
			ElectorMSP:   "Org2MSP",
		}

		It("Correct payload hash", func() {
			expected := []byte{151, 51, 62, 57, 99, 156, 115, 78, 195, 186, 169, 158, 44, 26, 175, 167, 211, 131, 185, 132, 86, 177, 37, 171, 42, 177, 205, 98, 55, 215, 151, 205}

			Expect(signature.HashElectorPayload()).Should(Equal(expected))
		})
	})

	Context("Unique key", func() {
		It("Correct unique key", func() {
			signature := &models.Signature{
				ElectionName:  "Best Crypto Currency",
				SignedMessage: CORRECT_SIGNATURE,
			}

			Expect(signature.UniqueKey()).Should(Equal("signature_Best Crypto Currency_89be4b7940344c5d1febebc5fdf73b100c749c025f7c3bc8bfa880d5852017b2"))
		})
	})
})
