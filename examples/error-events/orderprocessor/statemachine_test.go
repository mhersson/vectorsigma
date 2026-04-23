package orderprocessor

import (
	"testing"
)

func TestOrderProcessor_SuccessfulOrder(t *testing.T) {
	fsm := New()

	err := fsm.Run()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if fsm.ExtendedState.Data["fulfilled"] != true {
		t.Error("expected order to be fulfilled")
	}

	if fsm.ExtendedState.Data["error_handled"] == true {
		t.Error("error handler should not have been called")
	}
}

func TestOrderProcessor_PaymentFailureTriggersErrorEvent(t *testing.T) {
	fsm := New()

	// Simulate a payment failure
	fsm.ExtendedState.Data = make(map[string]interface{})
	fsm.ExtendedState.Data["simulate_payment_error"] = true

	err := fsm.Run()

	// The error should be stored in ExtendedState but Run() should complete
	if err == nil {
		t.Error("expected error to be returned")
	}

	// Verify error was handled
	if fsm.ExtendedState.Data["error_handled"] != true {
		t.Error("error handler should have been called via error event")
	}

	// Note: In composite states, all actions run even if one fails.
	// The error is only propagated after the composite state completes.
	// Since FulfillOrder runs after ChargePayment and also fails (payment not charged),
	// the final error will be from FulfillOrder, not ChargePayment.

	// Verify the error was captured
	if fsm.ExtendedState.Error == nil {
		t.Error("expected error to be stored in ExtendedState")
	}
}

func TestOrderProcessor_InitializationErrorTriggersErrorEvent(t *testing.T) {
	fsm := New()

	// Force an initialization error by clearing the data map after initialization
	// This will cause LoadSubject to fail
	originalInitAction := fsm.StateConfigs[Initializing].Composite.StateConfigs[InitializingContext].Actions[0].Execute

	fsm.StateConfigs[Initializing].Composite.StateConfigs[InitializingContext].Actions[0].Execute = func(params ...string) error {
		// Don't initialize the data map
		return nil
	}

	// Now LoadSubject will fail because it checks for initialized flag
	err := fsm.Run()

	// Restore original action
	fsm.StateConfigs[Initializing].Composite.StateConfigs[InitializingContext].Actions[0].Execute = originalInitAction

	if err == nil {
		t.Error("expected error to be returned")
	}

	// Verify error was handled
	if fsm.ExtendedState.Data["error_handled"] != true {
		t.Error("error handler should have been called via error event")
	}

	// Verify we didn't reach processing
	if fsm.ExtendedState.Data["validated"] == true {
		t.Error("should not have reached validation after initialization error")
	}
}

// BenchmarkOrderProcessor_SuccessfulFlow benchmarks the successful order flow.
func BenchmarkOrderProcessor_SuccessfulFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fsm := New()
		_ = fsm.Run()
	}
}

// BenchmarkOrderProcessor_ErrorHandling benchmarks the error handling flow.
func BenchmarkOrderProcessor_ErrorHandling(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fsm := New()
		fsm.ExtendedState.Data = make(map[string]interface{})
		fsm.ExtendedState.Data["simulate_payment_error"] = true
		_ = fsm.Run()
	}
}
