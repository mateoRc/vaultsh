package shell

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
)

type session struct {
	mu     sync.Mutex
	engine *Engine
}

type SessionManager struct {
	mu       sync.Mutex
	root     *filesystem.Directory
	sessions map[string]*session
}

func NewSessionManager(root *filesystem.Directory) *SessionManager {
	return &SessionManager{
		root:     root,
		sessions: make(map[string]*session),
	}
}

func (m *SessionManager) Execute(sessionID, line string) (command.Result, string, error) {
	current, sessionID, err := m.get(sessionID)
	if err != nil {
		return command.Result{}, "", err
	}

	current.mu.Lock()
	defer current.mu.Unlock()

	return current.engine.Execute(line), sessionID, nil
}

func (m *SessionManager) get(sessionID string) (*session, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if current, found := m.sessions[sessionID]; found {
		return current, sessionID, nil
	}

	newID, err := newSessionID()
	if err != nil {
		return nil, "", err
	}

	current := &session{engine: NewWithRoot(m.root)}
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
