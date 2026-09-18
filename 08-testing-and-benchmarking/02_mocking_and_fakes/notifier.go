package notifier

import (
	"errors"
	"fmt"
	"strings"
)

// MessageSender defines the interface dependency.
type MessageSender interface {
	SendMessage(to, body string) error
}

type OrderService struct {
	sender MessageSender
}

func NewOrderService(sender MessageSender) *OrderService {
	return &OrderService{sender: sender}
}

func (s *OrderService) ConfirmOrder(customerEmail string, orderID int) error {
	if !strings.Contains(customerEmail, "@") {
		return errors.New("invalid email address")
	}

	body := fmt.Sprintf("Your order #%d has been confirmed!", orderID)
	return s.sender.SendMessage(customerEmail, body)
}
