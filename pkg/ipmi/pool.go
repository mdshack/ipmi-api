package ipmi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bougou/go-ipmi"
)

// ClientPool manages a pool of IPMI clients
type ClientPool struct {
	connections map[string]*ipmi.Client
	mutex       sync.RWMutex
	poolSize    int
	timeout     time.Duration
}

// NewClientPool creates a new client pool
func NewClientPool(poolSize int, timeout time.Duration) *ClientPool {
	return &ClientPool{
		connections: make(map[string]*ipmi.Client),
		poolSize:    poolSize,
		timeout:     timeout,
	}
}

// GetClient retrieves or creates an IPMI client for the given credentials
func (cp *ClientPool) GetClient(host string, port int, username, password string) (*ipmi.Client, error) {
	key := fmt.Sprintf("%s:%d:%s", host, port, username)
	
	cp.mutex.RLock()
	client, exists := cp.connections[key]
	cp.mutex.RUnlock()
	
	if exists {
		return client, nil
	}
	
	// Create new client
	client, err := ipmi.NewClient(host, port, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPMI client: %w", err)
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cp.timeout)
	defer cancel()
	
	// Connect to the client
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to IPMI client: %w", err)
	}
	
	// Add to pool
	cp.mutex.Lock()
	// Check if we need to remove an existing connection to stay within pool size
	if len(cp.connections) >= cp.poolSize {
		// Remove the first connection we find
		for k, v := range cp.connections {
			// Create context for closing
			closeCtx, closeCancel := context.WithTimeout(context.Background(), cp.timeout)
			v.Close(closeCtx)
			closeCancel()
			delete(cp.connections, k)
			break
		}
	}
	cp.connections[key] = client
	cp.mutex.Unlock()
	
	return client, nil
}

// CloseAll closes all connections in the pool
func (cp *ClientPool) CloseAll() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	
	for _, client := range cp.connections {
		// Create context for closing
		ctx, cancel := context.WithTimeout(context.Background(), cp.timeout)
		client.Close(ctx)
		cancel()
	}
	cp.connections = make(map[string]*ipmi.Client)
}