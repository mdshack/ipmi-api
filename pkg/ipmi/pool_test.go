package ipmi

import (
	"testing"
	"time"
)

func TestClientPoolCreation(t *testing.T) {
	// Test creating a new client pool
	pool := NewClientPool(10, 30*time.Second)
	
	if pool == nil {
		t.Error("NewClientPool returned nil")
	}
	
	if pool.poolSize != 10 {
		t.Errorf("Expected poolSize to be 10, got %d", pool.poolSize)
	}
}