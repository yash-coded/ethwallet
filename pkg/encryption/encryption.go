package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strings"
)

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func EncryptData(data []byte, password string) (string, error) {
	key := deriveKey(password)

	block, err := aes.NewCipher(key)

	if err != nil {
		log.Fatalf("Failed to create new cipher: %v", err)
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("Failed to create new GCM: %v", err)
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())

	cipherText := aesGCM.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptData(encryptedData string, password string) (string, error) {
	key := deriveKey(password)

	block, err := aes.NewCipher(key)

	if err != nil {
		log.Fatalf("Failed to create new cipher: %v", err)
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("Failed to create new GCM: %v", err)
		return "", err
	}

	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)

	if err != nil {
		log.Fatalf("Failed to decode base64: %v", err)
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(decodedData) < nonceSize {
		log.Fatalf("Invalid data length")
		return "", err
	}

	nonce, cipherText := decodedData[:nonceSize], decodedData[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)

	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
		return "", err
	}

	return string(plainText), nil
}

func ConvertHexAddress(address string) string {
	address = strings.TrimPrefix(address, "0x")
	return strings.ToLower(address)
}
