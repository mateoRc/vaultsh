package shell

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/mateom/vaultsh/internal/command"
	"github.com/mateom/vaultsh/internal/filesystem"
	"github.com/mateom/vaultsh/internal/parser"
)

const (
	DefaultSessionTTL             = 30 * time.Minute
	DefaultSessionCleanupInterval = 5 * time.Minute
	DefaultMaxSessions            = 5000
)

var ErrSessionLimit = errors.New("session limit reached")

type SessionConfig struct {
	TTL         time.Duration
	Now         func() time.Time
	MaxSessions int
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
	max      int
	now      func() time.Time
	deps     Dependencies
}

func NewSessionManager(root *filesystem.Directory) *SessionManager {
	return newSessionManager(root, SessionConfig{}, Dependencies{})
}

func NewSessionManagerWithDependencies(
	root *filesystem.Directory,
	dependencies Dependencies,
) *SessionManager {
	return newSessionManager(root, SessionConfig{}, dependencies)
}

func NewSessionManagerWithConfig(
	root *filesystem.Directory,
	config SessionConfig,
) *SessionManager {
	return newSessionManager(root, config, Dependencies{})
}

func NewSessionManagerWithConfigAndDependencies(
	root *filesystem.Directory,
	config SessionConfig,
	dependencies Dependencies,
) *SessionManager {
	return newSessionManager(root, config, dependencies)
}

func newSessionManager(
	root *filesystem.Directory,
	config SessionConfig,
	dependencies Dependencies,
) *SessionManager {
	if config.TTL <= 0 {
		config.TTL = DefaultSessionTTL
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.MaxSessions <= 0 {
		config.MaxSessions = DefaultMaxSessions
	}

	return &SessionManager{
		root:     root,
		sessions: make(map[string]*session),
		ttl:      config.TTL,
		max:      config.MaxSessions,
		now:      config.Now,
		deps:     dependencies,
	}
}

func (m *SessionManager) Execute(
	sessionID,
	line string,
) (
	result command.Result,
	returnedSessionID string,
	currentDirectory string,
	err error,
) {
	current, sessionID, err := m.get(sessionID)
	if err != nil {
		return command.Result{}, "", "", err
	}

	current.mu.Lock()
	defer current.mu.Unlock()
	defer func() {
		current.lastActive = m.now()
	}()

	started := time.Now()
	defer func() {
		if recover() == nil {
			return
		}
		result = command.Result{
			Output:   "internal error",
			ExitCode: command.ExitFailure,
		}
		returnedSessionID = sessionID
		currentDirectory = current.engine.context.WorkingDirectory().Path()
		err = nil
		if name := commandName(line); name != "" && m.deps.Events != nil {
			_ = m.deps.Events.Record(
				"vault",
				"command.runtime_error",
				name,
				time.Since(started).Milliseconds(),
				result.ExitCode,
			)
		}
	}()

	result = current.engine.Execute(line)
	if name := commandName(line); name != "" && m.deps.Events != nil {
		_ = m.deps.Events.Record(
			"vault",
			"command.executed",
			name,
			time.Since(started).Milliseconds(),
			result.ExitCode,
		)
	}
	return result, sessionID, current.engine.context.WorkingDirectory().Path(), nil
}

func (m *SessionManager) Complete(
	sessionID string,
	line string,
	cursor int,
) (Completion, string, string, error) {
	current, sessionID, err := m.get(sessionID)
	if err != nil {
		return Completion{}, "", "", err
	}

	current.mu.Lock()
	defer current.mu.Unlock()
	defer func() {
		current.lastActive = m.now()
	}()

	return current.engine.Complete(line, cursor),
		sessionID,
		current.engine.context.WorkingDirectory().Path(),
		nil
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
	if len(m.sessions) >= m.max {
		return nil, "", ErrSessionLimit
	}

	newID, err := newSessionID()
	if err != nil {
		return nil, "", err
	}

	current := &session{
		engine: NewWithContextAndDependencies(
			NewExecutionContext(m.root),
			m.deps,
		),
		lastActive: m.now(),
	}
	m.sessions[newID] = current
	return current, newID, nil
}

func commandName(line string) string {
	tokens, err := parser.Tokenize(line)
	if err != nil || len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

func newSessionID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}
