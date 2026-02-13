package state

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x01}, masterKeySize)
	plain := []byte("ghp_example_token")

	ciphertext, err := encrypt(githubPATKey, plain, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	got, err := decrypt(githubPATKey, ciphertext, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(got, plain) {
		t.Fatalf("decrypted value mismatch: got=%q want=%q", got, plain)
	}
}

func TestDecryptFailsOnAADMismatch(t *testing.T) {
	key := bytes.Repeat([]byte{0x02}, masterKeySize)
	plain := []byte("ghp_example_token")

	ciphertext, err := encrypt(githubPATKey, plain, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if _, err := decrypt("another_secret", ciphertext, key); err == nil {
		t.Fatal("expected decrypt error for mismatched secret name")
	}
}

func TestLoadOrCreateMasterKey(t *testing.T) {
	tmp := t.TempDir()
	keyPath := filepath.Join(tmp, "master.key")

	key1, err := loadOrCreateMasterKey(keyPath)
	if err != nil {
		t.Fatalf("first load/create failed: %v", err)
	}
	if len(key1) != masterKeySize {
		t.Fatalf("unexpected key size: got %d want %d", len(key1), masterKeySize)
	}

	key2, err := loadOrCreateMasterKey(keyPath)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}
	if !bytes.Equal(key1, key2) {
		t.Fatal("expected same key on second load")
	}

	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("stat key failed: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("unexpected key permissions: got %o want 600", perm)
	}
}

func TestParseKeyFileLegacy(t *testing.T) {
	raw := bytes.Repeat([]byte{0xAB}, masterKeySize)
	key, err := parseKeyFile(raw)
	if err != nil {
		t.Fatalf("parseKeyFile legacy failed: %v", err)
	}
	if !bytes.Equal(key, raw) {
		t.Fatal("expected legacy key to be returned as-is")
	}
}

func TestParseKeyFileV1(t *testing.T) {
	expected := bytes.Repeat([]byte{0xCD}, masterKeySize)
	raw := make([]byte, 0, 1+masterKeySize)
	raw = append(raw, keyFileVersionV1)
	raw = append(raw, expected...)

	key, err := parseKeyFile(raw)
	if err != nil {
		t.Fatalf("parseKeyFile v1 failed: %v", err)
	}
	if !bytes.Equal(key, expected) {
		t.Fatal("expected v1 key bytes without version prefix")
	}
}

func TestParseKeyFileBadVersion(t *testing.T) {
	raw := make([]byte, 1+masterKeySize)
	raw[0] = 0xFF
	if _, err := parseKeyFile(raw); err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestParseKeyFileInvalidSize(t *testing.T) {
	raw := []byte{0x01, 0x02, 0x03} // 3 bytes — neither 32 nor 33
	if _, err := parseKeyFile(raw); err == nil {
		t.Fatal("expected error for invalid size")
	}
}

func TestLoadOrCreateMasterKeyV1Format(t *testing.T) {
	tmp := t.TempDir()
	keyPath := filepath.Join(tmp, "master.key")

	key, err := loadOrCreateMasterKey(keyPath)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Verify file is v1 format (33 bytes).
	raw, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("read key file failed: %v", err)
	}
	if len(raw) != 1+masterKeySize {
		t.Fatalf("expected %d byte v1 key file, got %d", 1+masterKeySize, len(raw))
	}
	if raw[0] != keyFileVersionV1 {
		t.Fatalf("expected version byte 0x%02x, got 0x%02x", keyFileVersionV1, raw[0])
	}
	if !bytes.Equal(raw[1:], key) {
		t.Fatal("key content mismatch")
	}
}
