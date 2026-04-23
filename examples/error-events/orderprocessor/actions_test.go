package orderprocessor_test

import (
	"testing"

	"github.com/mhersson/vectorsigma/examples/error-events/orderprocessor"
)

// +vectorsigma:action:ChargePayment
func TestOrderProcessor_ChargePaymentAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.ChargePaymentAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.ChargePaymentAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:FulfillOrder
func TestOrderProcessor_FulfillOrderAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.FulfillOrderAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.FulfillOrderAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:HandleError
func TestOrderProcessor_HandleErrorAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.HandleErrorAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.HandleErrorAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:InitializeContext
func TestOrderProcessor_InitializeContextAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.InitializeContextAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.InitializeContextAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:LoadSubject
func TestOrderProcessor_LoadSubjectAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.LoadSubjectAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.LoadSubjectAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:ValidateOrder
func TestOrderProcessor_ValidateOrderAction(t *testing.T) {
	type fields struct {
		context       *orderprocessor.Context
		currentState  orderprocessor.StateName
		stateConfigs  map[orderprocessor.StateName]orderprocessor.StateConfig
		ExtendedState *orderprocessor.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fsm := &orderprocessor.OrderProcessor{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.ValidateOrderAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("OrderProcessor.ValidateOrderAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
