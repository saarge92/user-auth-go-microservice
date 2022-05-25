package services

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type EncryptService struct {
	gcm cipher.AEAD
}

func NewCryptoService(key []byte) (*EncryptService, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &EncryptService{
		gcm: gcm,
	}, nil
}

func (s *EncryptService) Decrypt(cipherTextParam []byte) (plaintext []byte, err error) {
	cipherSource := make([]byte, base64.StdEncoding.DecodedLen(len(cipherTextParam)))
	_, err = base64.StdEncoding.Decode(cipherSource, cipherTextParam)
	if err != nil {
		return nil, err
	}

	nonce := cipherSource[0:s.gcm.NonceSize()]
	ciphertext := cipherSource[s.gcm.NonceSize():]

	plain, err := s.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, err
}
