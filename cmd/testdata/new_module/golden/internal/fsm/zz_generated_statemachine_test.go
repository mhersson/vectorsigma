// This file is generated by VectorSigma. DO NOT EDIT.
package fsm_test

import (
	"new_module/internal/fsm"
	"testing"
)

func TestTrafficLight_Run(t *testing.T) {
	type fields struct {
		Context       *fsm.Context
		CurrentState  fsm.StateName
		ExtendedState *fsm.ExtendedState
		StateConfigs  map[fsm.StateName]fsm.StateConfig
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Run the machine"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := fsm.New()
			fsm.Run()
		})
	}
}
