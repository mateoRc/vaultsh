package shell

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

const (
	DefaultSessionTTL             = 30 * time.Minute
	DefaultSessionCleanupInterval = 5 * time.Minute
)

type SessionConfig struct {
	TTL time.Duration
	Now func() time.Time
}

type session struct {
	mu         sync.Mutex
	engine     *Engine
	lastActive time.Time
}

type SessionManager struct {
	mu       sync.Mutex
	root     *filesystem.Directory
	sessions map[string]*session
	ttl      time.Duration
	now      func() time.Time
}

func NewSessionManager(root *filesystem.Directory) *SessionManager {
	return NewSessionManagerWithConfig(root, SessionConfig{})
}

func NewSessionManagerWithConfig(
	root *filesystem.Directory,
	config SessionConfig,
) *SessionManager {
	if config.TTL <= 0 {
		config.TTL = DefaultSessionTTL
	}
	if config.Now == nil {
		config.Now = time.Now
	}

	return &SessionManager{
		root:     root,
		sessions: make(map[string]*session),
		ttl:      config.TTL,
		now:      config.Now,
	}
}

func (m *SessionManager) Execute(sessionID, line string) (command.Result, string, error) {
	current, sessionID, err := m.get(sessionID)
	if err != nil {
		return command.Result{}, "", err
	}

	current.mu.Lock()
	defer current.mu.Unlock()
	defer func() {
		current.lastActive = m.now()
	}()

	return current.engine.Execute(line), sessionID, nil
}

func (m *SessionManager) Complete(
	sessionID string,
	line string,
	cursor int,
) (Completion, string, error) {
	current, sessionID, err := m.get(sessionID)
	if err != nil {
		return Completion{}, "", err
	}

	current.mu.Lock()
	defer current.mu.Unlock()
	defer func() {
		current.lastActive = m.now()
	}()

	return current.engine.Complete(line, cursor), sessionID, nil
}

func (m *SessionManager) CleanupExpired() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	removed := 0
	for sessionID, current := range m.sessions {
		current.mu.Lock()
		expired := now.Sub(current.lastActive) >= m.ttl
		current.mu.Unlock()
		if expired {
			delete(m.sessions, sessionID)
			removed++
		}
	}
	return removed
}

func (m *SessionManager) RunCleanup(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = DefaultSessionCleanupInterval
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.CleanupExpired()
		case <-ctx.Done():
			return
		}
	}
}

func (m *SessionManager) get(sessionID string) (*session, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if current, found := m.sessions[sessionID]; found {
		current.mu.Lock()
		expired := m.now().Sub(current.lastActive) >= m.ttl
		current.mu.Unlock()
		if !expired {
			return current, sessionID, nil
		}
		delete(m.sessions, sessionID)
	}

	newID, err := newSessionID()
	if err != nil {
		return nil, "", err
	}

	current := &session{
		engine:     NewWithRoot(m.root),
		lastActive: m.now(),
	}
	m.sessions[newID] = current
	return current, newID, nil
}

func newSessionID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}
