// service_test.go
package custody

import (
	"andi-custodian/internal/store"
	"context"
	"testing"
	"time"

	"andi-custodian/internal/wallet"
)

// MockSigner for testing
type MockSigner struct {
	signFunc func(ctx context.Context, req wallet.SignRequest) ([]byte, error)
}

func (m *MockSigner) Sign(ctx context.Context, req wallet.SignRequest) ([]byte, error) {
	if m.signFunc != nil {
		return m.signFunc(ctx, req)
	}
	return []byte("mock-signature"), nil
}

// These are real-looking, valid-length, checksum-compliant addresses.
// 💡 For Bitcoin tests later, use a valid testnet Bech32 address like:
// tb1q4d750u3s88c6mt8732j2q6gsn23rwwey25xxnm
const (
	testEthFrom = "0x8E76C1897e55d208b2b5f45cDb43FD7d403a9a31"
	testEthTo   = "0x742d35Cc6634C0532925a3b844Bc9dbd8b5E8a18"
)

func TestService_Transfer_Idempotency(t *testing.T) {
	signer := &MockSigner{}
	service := newTestService(t, signer)

	req := &TransferRequest{
		ID:    "req-123",
		Chain: "ethereum-sepolia",
		From:  testEthFrom, // valid
		To:    testEthTo,   // valid
		Value: "1",
	}

	// First call
	res1, err := service.Transfer(context.Background(), req)
	if err != nil {
		t.Fatalf("First transfer failed: %v", err)
	}

	// Second call with same ID
	res2, err := service.Transfer(context.Background(), req)
	if err != nil {
		t.Fatalf("Second transfer failed: %v", err)
	}

	if res1.TxID != res2.TxID {
		t.Error("Idempotency failed: TxID mismatch")
	}
}

func TestService_Transfer_SigningError(t *testing.T) {
	signer := &MockSigner{
		signFunc: func(ctx context.Context, req wallet.SignRequest) ([]byte, error) {
			return nil, wallet.ErrSigningFailed
		},
	}
	service := newTestService(t, signer)

	req := &TransferRequest{
		ID:    "req-456",
		Chain: "ethereum-sepolia",
		From:  "0x123",
		To:    "0x456",
		Value: "1",
	}

	_, err := service.Transfer(context.Background(), req)
	if err == nil {
		t.Fatal("Expected signing error")
	}
}

func TestService_Transfer_Monitoring(t *testing.T) {
	signer := &MockSigner{}
	service := newTestService(t, signer)

	req := &TransferRequest{
		ID:    "req-monitor",
		Chain: "ethereum-sepolia",
		From:  testEthFrom, // valid
		To:    testEthTo,   // valid
		Value: "1",
	}

	res, err := service.Transfer(context.Background(), req)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}

	if res.Status != "pending" {
		t.Errorf("Initial status = %s, want 'pending'", res.Status)
	}

	// Wait for simulated finality
	time.Sleep(6 * time.Second)

	// Check idempotency store for updated status
	if existing, ok := service.idempotency.Load(req.ID); ok {
		finalRes := existing.(*store.TransferResult)
		if finalRes.Status != "confirmed" {
			t.Errorf("Final status = %s, want 'confirmed'", finalRes.Status)
		}
	} else {
		t.Error("Result not found in idempotency store")
	}
}

func newTestService(t *testing.T, signer wallet.Signer) *Service {
	return NewService(signer, store.NewInMemoryStore())
}
