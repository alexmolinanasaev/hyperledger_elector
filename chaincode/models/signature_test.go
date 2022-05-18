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

			Expect(signature.UniqueKey()).Should(Equal("signature_Best Crypto Currency_3045022100c2b82ac0fcee4ef7bf22ef020b50a6bdfae0354961f9da94b8b818f715b3dc3702207361fd82bff0f8717890bf92662b8771f082137507626b3f70dc66299193a2cd"))
		})
	})
})
