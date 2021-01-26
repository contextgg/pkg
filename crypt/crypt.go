package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type Crypt interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(encrypted []byte) ([]byte, error)
}

type crypt struct {
	key []byte
}

// Encrypt will AES-encrypt the given byte slice
func (c *crypt) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure, therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

// DecryptBytes will AES-decrypt the given byte slice
func (c *crypt) Decrypt(encrypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if byteLen := len(encrypted); byteLen < aes.BlockSize {
		return nil, fmt.Errorf("invalid cipher size %d, expected at least %d", byteLen, aes.BlockSize)
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	// XORKeyStream can work in-place if the two arguments are the same.
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(encrypted, encrypted)
	return encrypted, nil
}

func NewCrypt(key []byte) (Crypt, error) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf("Invalid KEY to set for CRYPT_KEEPER_KEY; must be 16, 24, or 32 bytes (got %d)", keyLen)
	}

	return &crypt{
		key: key,
	}, nil
}
