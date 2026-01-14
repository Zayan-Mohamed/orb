package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

// NoiseHandshake implements simplified Noise_XX pattern for mutual authentication
// This provides perfect forward secrecy and mutual authentication
type NoiseHandshake struct {
	localEphemeral  *X25519KeyPair
	remoteEphemeral *[32]byte
	presharedKey    []byte // Derived from passcode
	initiator       bool
	handshakeHash   []byte
}

// NewNoiseHandshake creates a new Noise handshake
func NewNoiseHandshake(presharedKey []byte, initiator bool) (*NoiseHandshake, error) {
	if len(presharedKey) != 32 {
		return nil, errors.New("preshared key must be 32 bytes")
	}

	localEph, err := GenerateX25519KeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key: %w", err)
	}

	nh := &NoiseHandshake{
		localEphemeral: localEph,
		presharedKey:   presharedKey,
		initiator:      initiator,
		handshakeHash:  make([]byte, 0),
	}

	// Initialize handshake hash
	nh.updateHash(presharedKey)

	return nh, nil
}

// CreateInitiatorMessage creates the first handshake message (initiator -> responder)
// Message format: ephemeral_public_key || encrypted_auth
func (nh *NoiseHandshake) CreateInitiatorMessage() ([]byte, error) {
	if !nh.initiator {
		return nil, errors.New("only initiator can create initiator message")
	}

	// Update hash with our ephemeral public key
	nh.updateHash(nh.localEphemeral.Public[:])

	// Create authentication proof using preshared key
	authData := nh.computeAuthProof()

	// Encrypt auth data with preshared key
	cipher, err := NewAEAD(nh.presharedKey)
	if err != nil {
		return nil, err
	}

	encryptedAuth, err := cipher.Encrypt(authData)
	if err != nil {
		return nil, err
	}

	// Message: ephemeral_public || encrypted_auth
	message := make([]byte, 0, 32+len(encryptedAuth))
	message = append(message, nh.localEphemeral.Public[:]...)
	message = append(message, encryptedAuth...)

	return message, nil
}

// ProcessInitiatorMessage processes the initiator's message (responder side)
func (nh *NoiseHandshake) ProcessInitiatorMessage(message []byte) error {
	if nh.initiator {
		return errors.New("initiator cannot process initiator message")
	}

	if len(message) < 32 {
		return errors.New("message too short")
	}

	// Extract remote ephemeral public key
	var remotePub [32]byte
	copy(remotePub[:], message[:32])
	nh.remoteEphemeral = &remotePub

	// Update hash
	nh.updateHash(remotePub[:])

	// Decrypt and verify auth
	cipher, err := NewAEAD(nh.presharedKey)
	if err != nil {
		return err
	}

	authData, err := cipher.Decrypt(message[32:])
	if err != nil {
		return ErrAuthFailed
	}

	// Verify auth proof
	if !nh.verifyAuthProof(authData) {
		return ErrAuthFailed
	}

	return nil
}

// CreateResponderMessage creates the response message (responder -> initiator)
func (nh *NoiseHandshake) CreateResponderMessage() ([]byte, error) {
	if nh.initiator {
		return nil, errors.New("initiator cannot create responder message")
	}

	if nh.remoteEphemeral == nil {
		return nil, errors.New("must process initiator message first")
	}

	// Update hash with our ephemeral public key
	nh.updateHash(nh.localEphemeral.Public[:])

	// Create authentication proof
	authData := nh.computeAuthProof()

	// Compute shared secret for encryption
	sharedSecret, err := ComputeSharedSecret(&nh.localEphemeral.Private, nh.remoteEphemeral)
	if err != nil {
		return nil, err
	}

	// Derive encryption key from shared secret and handshake hash
	encKey := nh.deriveKey(sharedSecret[:], []byte("responder"))

	cipher, err := NewAEAD(encKey)
	if err != nil {
		return nil, err
	}

	encryptedAuth, err := cipher.Encrypt(authData)
	if err != nil {
		return nil, err
	}

	// Message: ephemeral_public || encrypted_auth
	message := make([]byte, 0, 32+len(encryptedAuth))
	message = append(message, nh.localEphemeral.Public[:]...)
	message = append(message, encryptedAuth...)

	return message, nil
}

// ProcessResponderMessage processes the responder's message (initiator side)
func (nh *NoiseHandshake) ProcessResponderMessage(message []byte) error {
	if !nh.initiator {
		return errors.New("responder cannot process responder message")
	}

	if len(message) < 32 {
		return errors.New("message too short")
	}

	// Extract remote ephemeral public key
	var remotePub [32]byte
	copy(remotePub[:], message[:32])
	nh.remoteEphemeral = &remotePub

	// Update hash
	nh.updateHash(remotePub[:])

	// Compute shared secret
	sharedSecret, err := ComputeSharedSecret(&nh.localEphemeral.Private, nh.remoteEphemeral)
	if err != nil {
		return err
	}

	// Derive decryption key
	decKey := nh.deriveKey(sharedSecret[:], []byte("responder"))

	cipher, err := NewAEAD(decKey)
	if err != nil {
		return err
	}

	authData, err := cipher.Decrypt(message[32:])
	if err != nil {
		return ErrAuthFailed
	}

	// Verify auth proof
	if !nh.verifyAuthProof(authData) {
		return ErrAuthFailed
	}

	return nil
}

// DeriveTransportKeys derives the final encryption keys for the tunnel
func (nh *NoiseHandshake) DeriveTransportKeys() (sendKey, recvKey []byte, err error) {
	if nh.remoteEphemeral == nil {
		return nil, nil, errors.New("handshake not complete")
	}

	// Compute final shared secret
	sharedSecret, err := ComputeSharedSecret(&nh.localEphemeral.Private, nh.remoteEphemeral)
	if err != nil {
		return nil, nil, err
	}

	// Keys must be complementary between initiator and responder:
	// What initiator sends = what responder receives
	// What initiator receives = what responder sends
	if nh.initiator {
		sendKey = nh.deriveKey(sharedSecret[:], []byte("initiator_to_responder"))
		recvKey = nh.deriveKey(sharedSecret[:], []byte("responder_to_initiator"))
	} else {
		sendKey = nh.deriveKey(sharedSecret[:], []byte("responder_to_initiator"))
		recvKey = nh.deriveKey(sharedSecret[:], []byte("initiator_to_responder"))
	}

	return sendKey, recvKey, nil
}

// updateHash updates the handshake hash (transcript)
func (nh *NoiseHandshake) updateHash(data []byte) {
	h := sha256.New()
	h.Write(nh.handshakeHash)
	h.Write(data)
	nh.handshakeHash = h.Sum(nil)
}

// deriveKey derives a key using HKDF-like construction
func (nh *NoiseHandshake) deriveKey(secret, info []byte) []byte {
	h := sha256.New()
	h.Write(nh.handshakeHash)
	h.Write(secret)
	h.Write(info)
	key := h.Sum(nil)
	return key[:32] // Return 32 bytes for ChaCha20-Poly1305
}

// computeAuthProof creates an authentication proof
func (nh *NoiseHandshake) computeAuthProof() []byte {
	h := sha256.New()
	h.Write(nh.handshakeHash)
	h.Write(nh.presharedKey)
	// Add random challenge for uniqueness
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		// This should never fail with crypto/rand, but handle it safely
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	h.Write(challenge)
	proof := h.Sum(nil)

	// Include challenge so remote can verify
	result := make([]byte, 0, 32+32)
	result = append(result, challenge...)
	result = append(result, proof...)
	return result
}

// verifyAuthProof verifies an authentication proof
func (nh *NoiseHandshake) verifyAuthProof(authData []byte) bool {
	if len(authData) != 64 {
		return false
	}

	challenge := authData[:32]
	receivedProof := authData[32:]

	// Recompute expected proof
	h := sha256.New()
	h.Write(nh.handshakeHash)
	h.Write(nh.presharedKey)
	h.Write(challenge)
	expectedProof := h.Sum(nil)

	// Constant-time comparison
	return ConstantTimeCompare(receivedProof, expectedProof)
}

// Cleanup securely erases sensitive data
func (nh *NoiseHandshake) Cleanup() {
	Zeroize(nh.localEphemeral.Private[:])
	Zeroize(nh.presharedKey)
	Zeroize(nh.handshakeHash)
}
