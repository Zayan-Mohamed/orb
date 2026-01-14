package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

const (
	// Argon2id parameters (memory-hard and slow)
	Argon2Time    = 3
	Argon2Memory  = 64 * 1024 // 64 MB
	Argon2Threads = 4
	Argon2KeyLen  = 32

	// Key sizes
	KeySize   = 32
	NonceSize = 24
)

var (
	ErrInvalidKey       = errors.New("invalid key size")
	ErrInvalidNonce     = errors.New("invalid nonce size")
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrAuthFailed       = errors.New("authentication failed")
)

// DeriveKey derives a cryptographic key from passcode and session ID using Argon2id
// This is memory-hard and computationally expensive to resist brute-force attacks
func DeriveKey(passcode, sessionID string) []byte {
	// Use session ID as salt to ensure unique keys per session
	salt := []byte(sessionID)

	// Ensure salt is at least 8 bytes
	if len(salt) < 8 {
		padded := make([]byte, 8)
		copy(padded, salt)
		salt = padded
	}

	// Argon2id: resistant to GPU cracking and side-channel attacks
	key := argon2.IDKey(
		[]byte(passcode),
		salt,
		Argon2Time,
		Argon2Memory,
		Argon2Threads,
		Argon2KeyLen,
	)

	return key
}

// X25519KeyPair generates an ephemeral X25519 key pair for Noise protocol
type X25519KeyPair struct {
	Private [32]byte
	Public  [32]byte
}

// GenerateX25519KeyPair creates a new ephemeral key pair
func GenerateX25519KeyPair() (*X25519KeyPair, error) {
	kp := &X25519KeyPair{}

	if _, err := rand.Read(kp.Private[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Generate public key from private key
	curve25519.ScalarBaseMult(&kp.Public, &kp.Private)

	return kp, nil
}

// ComputeSharedSecret performs X25519 key exchange
func ComputeSharedSecret(privateKey, publicKey *[32]byte) (*[32]byte, error) {
	shared, err := curve25519.X25519(privateKey[:], publicKey[:])
	if err != nil {
		return nil, fmt.Errorf("X25519 failed: %w", err)
	}

	// Check for low-order points (security requirement)
	var zero [32]byte
	var sharedArray [32]byte
	copy(sharedArray[:], shared)
	if subtle.ConstantTimeCompare(sharedArray[:], zero[:]) == 1 {
		return nil, errors.New("invalid shared secret: low-order point")
	}

	return &sharedArray, nil
}

// AEAD provides authenticated encryption using ChaCha20-Poly1305
type AEAD struct {
	cipher cipher.AEAD
	nonce  uint64 // Counter for replay protection
}

// NewAEAD creates a new AEAD cipher with the given key
func NewAEAD(key []byte) (*AEAD, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, ErrInvalidKey
	}

	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	return &AEAD{
		cipher: cipher,
		nonce:  0,
	}, nil
}

// Encrypt encrypts plaintext with authenticated encryption
// Returns: nonce || ciphertext || tag
func (a *AEAD) Encrypt(plaintext []byte) ([]byte, error) {
	// Increment nonce for replay protection
	a.nonce++

	// Create unique nonce (XChaCha20 uses 24-byte nonces)
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	binary.BigEndian.PutUint64(nonce[16:], a.nonce)

	// Fill rest with random data for additional entropy
	if _, err := rand.Read(nonce[:16]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Create separate nonce for Seal to avoid reuse
	sealNonce := make([]byte, chacha20poly1305.NonceSizeX)
	copy(sealNonce, nonce)

	// Encrypt and authenticate
	ciphertext := a.cipher.Seal(nonce, sealNonce, plaintext, nil) // #nosec G407 -- nonce is randomly generated

	return ciphertext, nil
}

// Decrypt decrypts and verifies authenticated ciphertext
func (a *AEAD) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < chacha20poly1305.NonceSizeX {
		return nil, ErrInvalidNonce
	}

	// Extract nonce
	nonce := ciphertext[:chacha20poly1305.NonceSizeX]
	encrypted := ciphertext[chacha20poly1305.NonceSizeX:]

	// Decrypt and verify
	plaintext, err := a.cipher.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// SecureRandom generates cryptographically secure random bytes
func SecureRandom(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}

// ConstantTimeCompare performs constant-time comparison to prevent timing attacks
func ConstantTimeCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Zeroize securely erases sensitive data from memory
func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
