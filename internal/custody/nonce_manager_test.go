// nonce_manager_test.go
package custody

import (
	"sync"
	"testing"
)

func TestNonceManager_GetNext(t *testing.T) {
	nm := NewNonceManager()
	addr := "0x123"

	nonce1 := nm.GetNext(addr)
	nonce2 := nm.GetNext(addr)

	if nonce1 != 0 {
		t.Errorf("First nonce = %d, want 0", nonce1)
	}
	if nonce2 != 1 {
		t.Errorf("Second nonce = %d, want 1", nonce2)
	}
}

func TestNonceManager_Concurrent(t *testing.T) {
	nm := NewNonceManager()
	addr := "0x123"
	const workers = 10
	const incs = 100

	var wg sync.WaitGroup
	nonces := make([]uint64, workers*incs)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < incs; j++ {
				idx := worker*incs + j
				nonces[idx] = nm.GetNext(addr)
			}
		}(i)
	}
	wg.Wait()

	// Check all nonces are unique and in range [0, workers*incs)
	seen := make(map[uint64]bool)
	for _, n := range nonces {
		if n >= uint64(workers*incs) {
			t.Fatalf("Nonce %d out of range", n)
		}
		if seen[n] {
			t.Fatalf("Duplicate nonce %d", n)
		}
		seen[n] = true
	}
}

func TestNonceManager_Reset(t *testing.T) {
	nm := NewNonceManager()
	addr := "0x123"

	nm.GetNext(addr)  // nonce = 0, next = 1
	nm.Reset(addr, 5) // jump to 5
	next := nm.GetNext(addr)

	if next != 5 {
		t.Errorf("After reset, nonce = %d, want 5", next)
	}
}
