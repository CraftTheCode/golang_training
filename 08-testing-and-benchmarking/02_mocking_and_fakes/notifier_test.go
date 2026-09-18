package notifier

import (
	"errors"
	"testing"
)

// MockSender records calls and allows stubbing errors
type MockSender struct {
	LastTo      string
	LastBody    string
	ShouldError bool
}

func (m *MockSender) SendMessage(to, body string) error {
	m.LastTo = to
	m.LastBody = body
	if m.ShouldError {
		return errors.New("network failure")
	}
	return nil
}

func TestOrderService_ConfirmOrder(t *testing.T) {
	t.Run("Successfully sends confirmation email", func(t *testing.T) {
		mock := &MockSender{}
		svc := NewOrderService(mock)

		err := svc.ConfirmOrder("customer@example.com", 1001)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if mock.LastTo != "customer@example.com" {
			t.Errorf("expected to=%q, got %q", "customer@example.com", mock.LastTo)
		}
		if mock.LastBody != "Your order #1001 has been confirmed!" {
			t.Errorf("unexpected body: %q", mock.LastBody)
		}
	})

	t.Run("Returns error when sender fails", func(t *testing.T) {
		mock := &MockSender{ShouldError: true}
		svc := NewOrderService(mock)

		err := svc.ConfirmOrder("customer@example.com", 1002)
		if err == nil {
			t.Errorf("expected error from downstream sender, got nil")
		}
	})

	t.Run("Validation error prevents sending", func(t *testing.T) {
		mock := &MockSender{}
		svc := NewOrderService(mock)

		err := svc.ConfirmOrder("invalid-email", 1003)
		if err == nil {
			t.Errorf("expected email validation error, got nil")
		}
		if mock.LastTo != "" {
			t.Errorf("sender was invoked despite invalid email")
		}
	})
}
