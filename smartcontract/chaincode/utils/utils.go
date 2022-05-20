package utils

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const ADMIN_IDENTITY = "Org1MSP.eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT"

func IsAdmin(ctx contractapi.TransactionContextInterface) (bool, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("empty mspID: %s", err)
	}

	userID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("empty ID: %s", err)
	}

	return fmt.Sprintf("%s.%s", mspID, userID) == ADMIN_IDENTITY, nil
}
