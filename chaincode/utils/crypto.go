package utils

import "crypto/ecdsa"

func GetAdminPub() *ecdsa.PublicKey {
	// TODO: захардкодить и распарсить ключ
	return &ecdsa.PublicKey{}
}

func CheckSignature(pub *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	return ecdsa.VerifyASN1(pub, hash, sig)
}
