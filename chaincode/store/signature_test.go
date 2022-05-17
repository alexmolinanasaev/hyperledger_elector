package store_test

import (
	"elector/chaincode/models"
	"elector/chaincode/store"
	"elector/chaincode/utils"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
)

var _ = Describe("Signature store", func() {
	electorChaincode := shimtest.NewMockStub(`signer`, cpaper_extended.NewCC())

	signatureStore := store.GetSignatureStore(electorChaincode)

	pubKey, err := utils.ExtractPubKeyFromCert([]byte(ADMIN_PUB_KEY))
	if err != nil {
		Fail("cannot extract pub key")
	}

	signature := &models.Signature{
		ElectionName: "Best Crypto Currency",
		ElectorMSP:   "Org2MSP",
		SignerPubKey: pubKey,
	}

	Context("Put one", func() {
		It("Signed message validation fail", func() {
			signature.SignedMessage = WRONG_SIGNATURE

			Expect(signatureStore.PutOne(signature)).Should(MatchError("wrong signature"))
		})

		It("Success", func() {
			// signature = &models.Signature{
			// 	ElectionName:  "Best Crypto Currency",
			// 	ElectorMSP:    "Org2MSP",
			// 	SignerPubKey:  pubKey,
			// 	SignedMessage: CORRECT_SIGNATURE,
			// }
			fmt.Println(signature)

			electorChaincode.MockTransactionStart("save signature")
			Expect(signatureStore.PutOne(signature)).Should(MatchError("wrong signature"))
			electorChaincode.MockTransactionEnd("save signature")
		})

		// It("Already exist", func() {
		// 	// signature = &models.Signature{
		// 	// 	ElectionName:  "Best Crypto Currency",
		// 	// 	ElectorMSP:    "Org2MSP",
		// 	// 	SignerPubKey:  pubKey,
		// 	// 	SignedMessage: CORRECT_SIGNATURE,
		// 	// }

		// 	electorChaincode.MockTransactionStart("save signature")
		// 	Expect(signatureStore.PutOne(signature)).Should(MatchError("already exist"))
		// 	electorChaincode.MockTransactionEnd("save signature")
		// })
	})

	// Context("Get one", func() {
	// 	It("Success", func() {
	// 		electorChaincode.MockTransactionStart("save signature")
	// 		Expect(signatureStore.GetOneByKey(signature.UniqueKey())).Should(Succeed())
	// 		electorChaincode.MockTransactionEnd("save signature")
	// 	})

	// 	It("Success", func() {
	// 		electorChaincode.MockTransactionStart("save signature")
	// 		Expect(signatureStore.GetOneByKey(signature.UniqueKey())).Should(Succeed())
	// 		electorChaincode.MockTransactionEnd("save signature")
	// 	})
	// })
})
