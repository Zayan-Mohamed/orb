package tunnel

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Zayan-Mohamed/orb/internal/crypto"
	"github.com/Zayan-Mohamed/orb/pkg/protocol"
	"github.com/gorilla/websocket"
)

const (
	// Timeout constants
	handshakeReadTimeout  = 120 * time.Second // Increased for slow connections 
	handshakeWriteTimeout = 30 * time.Second
	dataReadTimeout       = 120 * time.Second // Increased for large file transfers 
	dataWriteTimeout      = 30 * time.Second
)

// Tunnel represents an encrypted tunnel between peers
type Tunnel struct {
	conn       *websocket.Conn
	sendCipher *crypto.AEAD
	recvCipher *crypto.AEAD
	sessionID  string
	mu         sync.Mutex
	closed     bool
}

// NewTunnel creates a new encrypted tunnel
func NewTunnel(relayURL, sessionID, passcode string, isInitiator bool) (*Tunnel, error) {
	// Derive key from passcode
	presharedKey := crypto.DeriveKey(passcode, sessionID)

	// Connect to relay
	endpoint := "share"
	if !isInitiator {
		endpoint = "connect"
	}

	u, err := url.Parse(relayURL)
	if err != nil {
		return nil, fmt.Errorf("invalid relay URL: %w", err)
	}

	// Convert http(s) to ws(s)
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = "/" + endpoint
	q := u.Query()
	q.Set("session", sessionID)
	u.RawQuery = q.Encode()

	// Dial WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to relay: %w", err)
	}

	tunnel := &Tunnel{
		conn:      conn,
		sessionID: sessionID,
	}

	// Perform Noise handshake
	if err := tunnel.performHandshake(presharedKey, isInitiator); err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			return nil, fmt.Errorf("handshake failed: %w (failed to close: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("handshake failed: %w", err)
	}

	return tunnel, nil
}

// performHandshake performs the Noise protocol handshake
func (t *Tunnel) performHandshake(presharedKey []byte, isInitiator bool) error {
	noise, err := crypto.NewNoiseHandshake(presharedKey, isInitiator)
	if err != nil {
		return err
	}
	defer noise.Cleanup()

	if isInitiator {
		if err := t.performInitiatorHandshake(noise); err != nil {
			return err
		}
	} else {
		if err := t.performResponderHandshake(noise); err != nil {
			return err
		}
	}

	return t.setupTransportKeys(noise)
}

func (t *Tunnel) performInitiatorHandshake(noise *crypto.NoiseHandshake) error {
	// Send initiator message
	msg, err := noise.CreateInitiatorMessage()
	if err != nil {
		return err
	}

	frame := &protocol.Frame{
		Type:    protocol.FrameTypeHandshake,
		Payload: msg,
	}

	if err := t.sendRawFrame(frame); err != nil {
		return err
	}

	// Receive responder message
	respFrame, err := t.recvRawFrame()
	if err != nil {
		return err
	}

	if respFrame.Type != protocol.FrameTypeHandshakeResp {
		return fmt.Errorf("unexpected frame type: %d", respFrame.Type)
	}

	return noise.ProcessResponderMessage(respFrame.Payload)
}

func (t *Tunnel) performResponderHandshake(noise *crypto.NoiseHandshake) error {
	// Receive initiator message
	initFrame, err := t.recvRawFrame()
	if err != nil {
		return err
	}

	if initFrame.Type != protocol.FrameTypeHandshake {
		return fmt.Errorf("unexpected frame type: %d", initFrame.Type)
	}

	if err := noise.ProcessInitiatorMessage(initFrame.Payload); err != nil {
		return err
	}

	// Send responder message
	msg, err := noise.CreateResponderMessage()
	if err != nil {
		return err
	}

	frame := &protocol.Frame{
		Type:    protocol.FrameTypeHandshakeResp,
		Payload: msg,
	}

	return t.sendRawFrame(frame)
}

func (t *Tunnel) setupTransportKeys(noise *crypto.NoiseHandshake) error {
	// Derive transport keys
	sendKey, recvKey, err := noise.DeriveTransportKeys()
	if err != nil {
		return err
	}

	// Create ciphers for secure transport
	t.sendCipher, err = crypto.NewAEAD(sendKey)
	if err != nil {
		return err
	}

	t.recvCipher, err = crypto.NewAEAD(recvKey)
	if err != nil {
		return err
	}

	// Cleanup keys from memory
	crypto.Zeroize(sendKey)
	crypto.Zeroize(recvKey)

	return nil
}

// SendFrame sends an encrypted frame
func (t *Tunnel) SendFrame(frame *protocol.Frame) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("tunnel closed")
	}

	// Serialize frame payload
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(frame); err != nil {
		return fmt.Errorf("failed to encode frame: %w", err)
	}

	// Encrypt payload
	encrypted, err := t.sendCipher.Encrypt(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	// Send over WebSocket
	_ = t.conn.SetWriteDeadline(time.Now().Add(dataWriteTimeout))
	if err := t.conn.WriteMessage(websocket.BinaryMessage, encrypted); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

// ReceiveFrame receives and decrypts a frame
func (t *Tunnel) ReceiveFrame() (*protocol.Frame, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil, fmt.Errorf("tunnel closed")
	}

	// Receive from WebSocket
	_ = t.conn.SetReadDeadline(time.Now().Add(dataReadTimeout))
	_, encrypted, err := t.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to receive: %w", err)
	}

	// Decrypt payload
	decrypted, err := t.recvCipher.Decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Deserialize frame
	var frame protocol.Frame
	dec := gob.NewDecoder(bytes.NewReader(decrypted))
	if err := dec.Decode(&frame); err != nil {
		return nil, fmt.Errorf("failed to decode frame: %w", err)
	}

	// Validate frame type
	if !protocol.ValidateFrameType(frame.Type) {
		return nil, protocol.ErrUnknownFrameType
	}

	return &frame, nil
}

// sendRawFrame sends an unencrypted frame (for handshake only)
func (t *Tunnel) sendRawFrame(frame *protocol.Frame) error {
	var buf bytes.Buffer
	if err := protocol.WriteFrame(&buf, frame); err != nil {
		return err
	}

	_ = t.conn.SetWriteDeadline(time.Now().Add(handshakeWriteTimeout))
	return t.conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

// recvRawFrame receives an unencrypted frame (for handshake only)
func (t *Tunnel) recvRawFrame() (*protocol.Frame, error) {
	_ = t.conn.SetReadDeadline(time.Now().Add(handshakeReadTimeout))
	_, data, err := t.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return protocol.ReadFrame(bytes.NewReader(data))
}

// Ping sends a ping and waits for pong
func (t *Tunnel) Ping() error {
	frame := &protocol.Frame{
		Type:    protocol.FrameTypePing,
		Payload: []byte{},
	}

	if err := t.SendFrame(frame); err != nil {
		return err
	}

	// Wait for pong
	resp, err := t.ReceiveFrame()
	if err != nil {
		return err
	}

	if resp.Type != protocol.FrameTypePong {
		return fmt.Errorf("expected pong, got %d", resp.Type)
	}

	return nil
}

// Close closes the tunnel
func (t *Tunnel) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true
	return t.conn.Close()
}

// IsClosed returns whether the tunnel is closed
func (t *Tunnel) IsClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closed
}
