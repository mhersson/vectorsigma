//go:build integration

package fsm_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"operator/output/fsm"
)

func init() {
	// Register custom setup hook
	CustomSetupHook = func() error {
		return nil
	}

	// Register custom teardown hook
	CustomTeardownHook = func() error {
		return nil
	}
}

func TestTestreconcileloop_Run(t *testing.T) {
	type fields struct {
		Context       *fsm.Context
		CurrentState  fsm.StateName
		ExtendedState *fsm.ExtendedState
		StateConfigs  map[fsm.StateName]fsm.StateConfig
	}
	tests := []struct {
		name    string
		want    ctrl.Result
		wantErr bool
	}{
		{name: "Happy path", want: ctrl.Result{}, wantErr: false},
	}
	for _, tt := range tests {
		setup(t)
		teardown(t)
		t.Run(tt.name, func(t *testing.T) {
			fsm := fsm.New()
			fsm.Context.Client = k8sClient
			fsm.Context.Ctx = context.TODO()
			fsm.Context.Logger = logr.Discard()
			fsm.ExtendedState.ResourceName = resourceName
			got, err := fsm.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Testreconcileloop.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Testreconcileloop.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
