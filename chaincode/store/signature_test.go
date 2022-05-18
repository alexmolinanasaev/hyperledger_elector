package store_test

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"elector/chaincode/utils"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
)

var _ = Describe("Signature store", func() {
	electorChaincode := shimtest.NewMockStub(`elector`, cpaper_extended.NewCC())

	signatureStore := store.GetSignatureStore(electorChaincode)

	pubKey, err := utils.ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
	if err != nil {
		Fail("cannot extract pub key")
	}

	signature := &models.Signature{
		ElectionName:  "Best Crypto Currency",
		ElectorMSP:    "Org2MSP",
		SignerPubKey:  pubKey,
		SignedMessage: WRONG_SIGNATURE,
	}

	Context("Put one", func() {
		It("Signed message validation fail", func() {
			Expect(signatureStore.PutOne(signature)).Should(MatchError("wrong signature"))
		})

		It("Success", func() {
			signature.SignedMessage = CORRECT_SIGNATURE

			electorChaincode.MockTransactionStart("save signature")
			Expect(signatureStore.PutOne(signature)).Should(Succeed())
			electorChaincode.MockTransactionEnd("save signature")
		})

		It("Already exist", func() {
			electorChaincode.MockTransactionStart("save signature")
			Expect(signatureStore.PutOne(signature)).Should(MatchError("already exist"))
			electorChaincode.MockTransactionEnd("save signature")
		})
	})

	Context("Get one", func() {
		It("Not found", func() {
			electorChaincode.MockTransactionStart("get signature")

			s, err := signatureStore.GetOneByKey("wrongKey")
			Expect(s).Should(BeNil())
			Expect(err).Should(BeNil())

			electorChaincode.MockTransactionEnd("get signature")
		})

		It("Success", func() {
			s := &models.Signature{
				ElectionName:  "Best Crypto Currency",
				SignedMessage: CORRECT_SIGNATURE,
			}

			electorChaincode.MockTransactionStart("get signature")
			Expect(signatureStore.GetOneByKey(signature.UniqueKey())).Should(Equal(s))
			electorChaincode.MockTransactionEnd("get signature")
		})
	})
})
