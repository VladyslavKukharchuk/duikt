package service

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type SignatureService struct{}

func NewSignatureService() *SignatureService {
	return &SignatureService{}
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(message))
	messageHash := hash.Sum(nil)

	signature, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func VerifySignature(message string, signature []byte, publicKey *ecdsa.PublicKey) (bool, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(message))
	messageHash := hash.Sum(nil)

	sigPublicKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return false, err
	}

	return publicKey.X.Cmp(sigPublicKey.X) == 0 && publicKey.Y.Cmp(sigPublicKey.Y) == 0, nil
}
