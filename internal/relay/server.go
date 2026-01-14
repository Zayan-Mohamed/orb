package relay

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Zayan-Mohamed/orb/internal/session"
	"github.com/gorilla/websocket"
)

const (
	// WebSocket settings
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2 * 1024 * 1024 // 2 MB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

// RelayServer is the blind relay server that forwards encrypted bytes
type RelayServer struct {
	sessionManager *session.SessionManager
	connections    map[string]*ConnectionPair
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// ConnectionPair represents a sharer-receiver connection pair
type ConnectionPair struct {
	SessionID string
	Sharer    *websocket.Conn
	Receiver  *websocket.Conn
	mu        sync.Mutex
	created   time.Time
	lastPing  time.Time
}

// NewRelayServer creates a new relay server
func NewRelayServer() *RelayServer {
	ctx, cancel := context.WithCancel(context.Background())

	rs := &RelayServer{
		sessionManager: session.NewSessionManager(),
		connections:    make(map[string]*ConnectionPair),
		ctx:            ctx,
		cancel:         cancel,
	}

	// Start connection monitor
	go rs.monitorConnections()

	return rs
}

// HandleShare handles the share endpoint (initiator)
func (rs *RelayServer) HandleShare(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "session required", http.StatusBadRequest)
		return
	}

	// Validate session exists
	sess, exists := rs.sessionManager.GetSession(sessionID)
	if !exists {
		http.Error(w, "invalid session", http.StatusNotFound)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Configure connection
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	rs.mu.Lock()
	pair, exists := rs.connections[sessionID]
	if !exists {
		pair = &ConnectionPair{
			SessionID: sessionID,
			Sharer:    conn,
			created:   time.Now(),
			lastPing:  time.Now(),
		}
		rs.connections[sessionID] = pair
	} else {
		pair.Sharer = conn
	}
	rs.mu.Unlock()

	log.Printf("Sharer connected: session=%s", sessionID)

	// Start message forwarding
	go rs.forwardMessages(conn, sessionID, true)
	go rs.keepAlive(conn)

	// Update session activity
	rs.sessionManager.UpdateActivity(sessionID)

	// Mark session as active
	sess.Active = true
}

// HandleConnect handles the connect endpoint (receiver)
func (rs *RelayServer) HandleConnect(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "session required", http.StatusBadRequest)
		return
	}

	// Validate session
	_, exists := rs.sessionManager.GetSession(sessionID)
	if !exists {
		http.Error(w, "invalid session", http.StatusNotFound)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Configure connection
	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	rs.mu.Lock()
	pair, exists := rs.connections[sessionID]
	if !exists {
		pair = &ConnectionPair{
			SessionID: sessionID,
			Receiver:  conn,
			created:   time.Now(),
			lastPing:  time.Now(),
		}
		rs.connections[sessionID] = pair
	} else {
		pair.Receiver = conn
	}
	rs.mu.Unlock()

	log.Printf("Receiver connected: session=%s", sessionID)

	// Start message forwarding
	go rs.forwardMessages(conn, sessionID, false)
	go rs.keepAlive(conn)

	// Update session activity
	rs.sessionManager.UpdateActivity(sessionID)
}

// forwardMessages forwards encrypted messages between peers
// The relay server never sees plaintext - it's a blind pipe
func (rs *RelayServer) forwardMessages(conn *websocket.Conn, sessionID string, isSharer bool) {
	defer func() {
		conn.Close()
		rs.cleanupConnection(sessionID, isSharer)
	}()

	for {
		// Read encrypted message (the relay is blind to content)
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Never log the message content (privacy requirement)

		// Forward to the other peer
		rs.mu.RLock()
		pair, exists := rs.connections[sessionID]
		rs.mu.RUnlock()

		if !exists {
			break
		}

		pair.mu.Lock()
		var target *websocket.Conn
		if isSharer && pair.Receiver != nil {
			target = pair.Receiver
		} else if !isSharer && pair.Sharer != nil {
			target = pair.Sharer
		}

		if target != nil {
			_ = target.SetWriteDeadline(time.Now().Add(writeWait))
			if err := target.WriteMessage(messageType, message); err != nil {
				log.Printf("Failed to forward message: %v", err)
				pair.mu.Unlock()
				break
			}
		}
		pair.mu.Unlock()

		// Update activity
		rs.sessionManager.UpdateActivity(sessionID)
	}
}

// keepAlive sends periodic pings to keep connection alive
func (rs *RelayServer) keepAlive(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-rs.ctx.Done():
			return
		}
	}
}

// cleanupConnection removes a connection from the pair
func (rs *RelayServer) cleanupConnection(sessionID string, isSharer bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	pair, exists := rs.connections[sessionID]
	if !exists {
		return
	}

	if isSharer {
		pair.Sharer = nil
	} else {
		pair.Receiver = nil
	}

	// If both connections are gone, remove the pair
	if pair.Sharer == nil && pair.Receiver == nil {
		delete(rs.connections, sessionID)
		log.Printf("Session closed: %s", sessionID)
	}
}

// monitorConnections monitors and cleans up stale connections
func (rs *RelayServer) monitorConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rs.mu.Lock()
			now := time.Now()
			for sessionID, pair := range rs.connections {
				// Remove stale connections (30 minutes inactive)
				if now.Sub(pair.lastPing) > 30*time.Minute {
					if pair.Sharer != nil {
						pair.Sharer.Close()
					}
					if pair.Receiver != nil {
						pair.Receiver.Close()
					}
					delete(rs.connections, sessionID)
					log.Printf("Removed stale connection: %s", sessionID)
				}
			}
			rs.mu.Unlock()
		case <-rs.ctx.Done():
			return
		}
	}
}

// HandleCreateSession handles session creation
func (rs *RelayServer) HandleCreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SharedPath string `json:"shared_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Create session
	sess, err := rs.sessionManager.CreateSession(req.SharedPath)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	// Return session details
	response := map[string]string{
		"session_id": sess.ID,
		"passcode":   sess.Passcode,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)

	// Never log passcodes (security requirement)
	log.Printf("Session created: %s", sess.ID)
}

// Start starts the relay server
func (rs *RelayServer) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/share", rs.HandleShare)
	mux.HandleFunc("/connect", rs.HandleConnect)
	mux.HandleFunc("/session/create", rs.HandleCreateSession)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Relay server starting on %s", addr)
	return server.ListenAndServe()
}

// Shutdown gracefully shuts down the relay server
func (rs *RelayServer) Shutdown() {
	rs.cancel()

	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Close all connections
	for _, pair := range rs.connections {
		if pair.Sharer != nil {
			pair.Sharer.Close()
		}
		if pair.Receiver != nil {
			pair.Receiver.Close()
		}
	}

	rs.connections = make(map[string]*ConnectionPair)
}
