package orderprocessor

import (
	"errors"
	"fmt"
)

// InitializeContext initializes the order context.
func (fsm *OrderProcessor) InitializeContextAction(_ ...string) error {
	fsm.Context.Logger.Info("initializing context")

	// Simulate initialization logic
	if fsm.ExtendedState.Data == nil {
		fsm.ExtendedState.Data = make(map[string]interface{})
	}

	fsm.ExtendedState.Data["initialized"] = true

	return nil
}

// LoadSubject loads order subject information.
func (fsm *OrderProcessor) LoadSubjectAction(_ ...string) error {
	fsm.Context.Logger.Info("loading subject")

	// Check if context was initialized
	if fsm.ExtendedState.Data["initialized"] != true {
		return errors.New("context not initialized")
	}

	fsm.ExtendedState.Data["subject"] = "Order #12345"

	return nil
}

// ValidateOrder validates the order.
func (fsm *OrderProcessor) ValidateOrderAction(_ ...string) error {
	fsm.Context.Logger.Info("validating order")

	subject, ok := fsm.ExtendedState.Data["subject"].(string)
	if !ok {
		return errors.New("no subject loaded")
	}

	fsm.Context.Logger.Info("order validated", "subject", subject)
	fsm.ExtendedState.Data["validated"] = true

	return nil
}

// ChargePayment processes payment.
func (fsm *OrderProcessor) ChargePaymentAction(_ ...string) error {
	fsm.Context.Logger.Info("charging payment")

	if fsm.ExtendedState.Data["validated"] != true {
		return errors.New("order not validated")
	}

	// Simulate a payment failure to test error handling
	if fsm.ExtendedState.Data["simulate_payment_error"] == true {
		return errors.New("payment declined")
	}

	fsm.ExtendedState.Data["payment_charged"] = true

	return nil
}

// FulfillOrder fulfills the order.
func (fsm *OrderProcessor) FulfillOrderAction(_ ...string) error {
	fsm.Context.Logger.Info("fulfilling order")

	if fsm.ExtendedState.Data["payment_charged"] != true {
		return errors.New("payment not charged")
	}

	fsm.ExtendedState.Data["fulfilled"] = true

	return nil
}

// HandleError handles errors.
func (fsm *OrderProcessor) HandleErrorAction(_ ...string) error {
	fsm.Context.Logger.Error("handling error", "error", fsm.ExtendedState.Error)

	// Initialize Data map if needed
	if fsm.ExtendedState.Data == nil {
		fsm.ExtendedState.Data = make(map[string]interface{})
	}

	// Log error details
	if fsm.ExtendedState.Error != nil {
		fmt.Printf("Error occurred: %v\n", fsm.ExtendedState.Error)
	}

	// Mark error as handled
	fsm.ExtendedState.Data["error_handled"] = true

	return nil
}
