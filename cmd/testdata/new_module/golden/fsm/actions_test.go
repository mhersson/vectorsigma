package fsm_test

import (
	"testing"

	"new_module/fsm"
)

func TestTrafficLight_SwitchInAction(t *testing.T) {
	type fields struct {
		context       *fsm.Context
		currentState  fsm.StateName
		stateConfigs  map[fsm.StateName]fsm.StateConfig
		ExtendedState *fsm.ExtendedState
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
			fsm := &fsm.TrafficLight{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.SwitchInAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("TrafficLight.SwitchInAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
