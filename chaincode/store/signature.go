package store

import (
	"elector/chaincode/models"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type SignatureStore struct {
	store *store
}

func GetSignatureStore(stub shim.ChaincodeStubInterface) *SignatureStore {
	return &SignatureStore{
		store: getStore(stub),
	}
}

func (s *SignatureStore) PutOne(signature *models.Signature) error {
	// сигнатура должна пройти 2 стадии валидации
	// 1) проверка самого содержимого
	// 2) проверка хэша
	// это связано с тем, что перед сохранением лишние данные, которые используются для валидации удаляются
	// перед непосредственных сохранением валидация впринципе не имеет смысла, но она нужна для удовлетворения
	// алгоритмам общего хранения
	if err := signature.Validate(); err != nil {
		return err
	}

	foundSignature, err := s.GetOneByKey(signature.UniqueKey())
	if err != nil {
		return fmt.Errorf("cannot verify if already exist: %s", err)
	}
	// fmt.Println(foundSignature)

	if foundSignature != nil {
		return fmt.Errorf("already exist")
	}

	messageHash := &models.Signature{
		MessageHash: signature.HashElectorPayload(),
	}

	return s.store.putOne(messageHash)
}

func (s *SignatureStore) PutMany(signature []*models.Signature) error {
	for _, d := range signature {
		err := s.PutOne(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SignatureStore) GetOneByKey(key string) (*models.Signature, error) {
	signatureRaw, err := s.store.getOneByKey(key)
	if err != nil {
		return nil, fmt.Errorf("cannot get signature by key: %s", err)
	}

	if signatureRaw == nil {
		return nil, nil
	}

	result := &models.Signature{}
	if err := json.Unmarshal(signatureRaw, result); err != nil {
		return nil, fmt.Errorf("cannot unmarshal signature: %s", err)
	}

	return result, nil
}
