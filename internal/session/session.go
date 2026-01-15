package session

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

const (
	SessionIDLength   = 6 // e.g., "7F9Q2A"
	PasscodeFormat    = 3 // e.g., "493-771"
	SessionTimeout    = 24 * time.Hour
	MaxFailedAttempts = 5
)

// Session represents an active tunnel session
type Session struct {
	ID             string
	Passcode       string
	Created        time.Time
	LastActivity   time.Time
	FailedAttempts int
	Locked         bool
	SharedPath     string
	Active         bool
	ConnectedPeer  string
}

// SessionManager manages all active sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}

	// Start cleanup goroutine
	go sm.cleanupExpired()

	return sm
}

// GenerateSessionID creates a random, human-readable session ID
func GenerateSessionID() (string, error) {
	// Use crypto/rand for security
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Encode to base32 (no ambiguous characters)
	encoded := base32.StdEncoding.EncodeToString(bytes)
	// Remove padding and take first 6 chars
	sessionID := strings.TrimRight(encoded, "=")[:SessionIDLength]

	return sessionID, nil
}

// GeneratePasscode creates a random numeric passcode in format "XXX-XXX"
func GeneratePasscode() (string, error) {
	// Generate 6-digit number (000000-999999)
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate passcode: %w", err)
	}

	// Format as XXX-XXX
	code := fmt.Sprintf("%06d", n.Int64())
	passcode := fmt.Sprintf("%s-%s", code[:3], code[3:])

	return passcode, nil
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(sharedPath string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Generate unique session ID
	var sessionID string
	var err error
	for {
		sessionID, err = GenerateSessionID()
		if err != nil {
			return nil, err
		}

		// Ensure uniqueness
		if _, exists := sm.sessions[sessionID]; !exists {
			break
		}
	}

	passcode, err := GeneratePasscode()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:           sessionID,
		Passcode:     passcode,
		Created:      time.Now(),
		LastActivity: time.Now(),
		SharedPath:   sharedPath,
		Active:       true,
	}

	sm.sessions[sessionID] = session

	return session, nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	return session, exists
}

// ValidatePasscode validates a passcode for a session (with rate limiting)
func (sm *SessionManager) ValidatePasscode(sessionID, passcode string) error {
	// Start timer for constant-time response
	start := time.Now()
	defer func() {
		// Ensure function always takes exactly 100ms to mitigate timing attacks
		elapsed := time.Since(start)
		remaining := (100 * time.Millisecond) - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		// Return generic error to prevent enumeration
		return fmt.Errorf("authentication failed")
	}

	// Check if locked
	if session.Locked {
		return fmt.Errorf("session locked due to too many failed attempts")
	}

	// Check if expired
	if time.Since(session.Created) > SessionTimeout {
		delete(sm.sessions, sessionID)
		return fmt.Errorf("session expired")
	}

	// Validate passcode (constant-time comparison)
	if !constantTimeStringCompare(session.Passcode, passcode) {
		session.FailedAttempts++
		if session.FailedAttempts >= MaxFailedAttempts {
			session.Locked = true
			return fmt.Errorf("session locked due to too many failed attempts")
		}
		return fmt.Errorf("authentication failed")
	}

	// Success - reset failed attempts
	session.FailedAttempts = 0
	session.LastActivity = time.Now()

	return nil
}

// UpdateActivity updates the last activity timestamp
func (sm *SessionManager) UpdateActivity(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		session.LastActivity = time.Now()
	}
}

// RevokeSession terminates a session
func (sm *SessionManager) RevokeSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.Active = false
	delete(sm.sessions, sessionID)

	return nil
}

// ListSessions returns all active sessions
func (sm *SessionManager) ListSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		if session.Active {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// cleanupExpired removes expired sessions periodically
func (sm *SessionManager) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for id, session := range sm.sessions {
			// Remove sessions that are expired or inactive for too long
			if now.Sub(session.Created) > SessionTimeout ||
				now.Sub(session.LastActivity) > 30*time.Minute {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}

// constantTimeStringCompare performs constant-time string comparison
func constantTimeStringCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i] ^ b[i])
	}

	return result == 0
}
