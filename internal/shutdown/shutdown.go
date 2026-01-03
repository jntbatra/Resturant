package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Manager manages graceful shutdown of application
type Manager struct {
	shutdownTimeout time.Duration
	hooks           []ShutdownHook
	mu              sync.Mutex
	shutdownOnce    sync.Once
	shutdownChan    chan struct{}
}

// ShutdownHook is a function called during graceful shutdown
type ShutdownHook func(ctx context.Context) error

// NewManager creates a new shutdown manager
func NewManager(timeout time.Duration) *Manager {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Manager{
		shutdownTimeout: timeout,
		hooks:           make([]ShutdownHook, 0),
		shutdownChan:    make(chan struct{}),
	}
}

// RegisterHook registers a shutdown hook to be called during graceful shutdown
func (m *Manager) RegisterHook(hook ShutdownHook) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hooks = append(m.hooks, hook)
}

// Wait blocks until shutdown signal is received
func (m *Manager) Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	m.Shutdown()
}

// Shutdown initiates graceful shutdown
func (m *Manager) Shutdown() {
	m.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
		defer cancel()

		m.executeShutdownHooks(ctx)
		close(m.shutdownChan)
	})
}

// executeShutdownHooks executes all registered shutdown hooks
func (m *Manager) executeShutdownHooks(ctx context.Context) {
	m.mu.Lock()
	hooks := make([]ShutdownHook, len(m.hooks))
	copy(hooks, m.hooks)
	m.mu.Unlock()

	// Execute hooks in reverse order (LIFO)
	for i := len(hooks) - 1; i >= 0; i-- {
		hook := hooks[i]
		if err := hook(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}

	log.Println("Graceful shutdown completed")
}

// Done returns a channel that closes when shutdown is complete
func (m *Manager) Done() <-chan struct{} {
	return m.shutdownChan
}

// IsShuttingDown checks if shutdown has been initiated
func (m *Manager) IsShuttingDown() bool {
	select {
	case <-m.shutdownChan:
		return true
	default:
		return false
	}
}
