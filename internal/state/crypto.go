package state

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	masterKeySize = 32
	nonceSize     = 12
)

var (
	errInvalidCiphertext = errors.New("invalid ciphertext")
)

// loadOrCreateMasterKey loads a 32-byte key from disk or creates one with 0600 permissions.
// Key file format v1: [version_byte=0x01][32_bytes_key] = 33 bytes total.
// Legacy format (v0): [32_bytes_key] = 32 bytes, read transparently.
func loadOrCreateMasterKey(path string) ([]byte, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create key dir: %w", err)
	}

	raw, err := os.ReadFile(path)
	if err == nil {
		return parseKeyFile(raw)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	key := make([]byte, masterKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	// Write key with version header.
	versioned := make([]byte, 0, 1+masterKeySize)
	versioned = append(versioned, keyFileVersionV1)
	versioned = append(versioned, key...)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return loadOrCreateMasterKey(path)
		}
		return nil, fmt.Errorf("create key file: %w", err)
	}

	if _, err := file.Write(versioned); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("write key file: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close key file: %w", err)
	}

	return key, nil
}

const keyFileVersionV1 = byte(0x01)

// parseKeyFile handles both legacy (32-byte) and versioned (33-byte) key files.
func parseKeyFile(raw []byte) ([]byte, error) {
	switch len(raw) {
	case masterKeySize: // Legacy v0: raw 32-byte key
		return raw, nil
	case 1 + masterKeySize: // v1: version byte + 32-byte key
		if raw[0] != keyFileVersionV1 {
			return nil, fmt.Errorf("unsupported key file version: 0x%02x", raw[0])
		}
		return raw[1:], nil
	default:
		return nil, fmt.Errorf("invalid key file: expected %d or %d bytes, got %d", masterKeySize, 1+masterKeySize, len(raw))
	}
}

func encrypt(secretName string, plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, []byte(secretName))
	out := make([]byte, 0, len(nonce)+len(ciphertext))
	out = append(out, nonce...)
	out = append(out, ciphertext...)
	return out, nil
}

func decrypt(secretName string, ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) <= nonceSize {
		return nil, errInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	nonce := ciphertext[:nonceSize]
	data := ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, data, []byte(secretName))
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return plaintext, nil
}
