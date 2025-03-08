package fsm_test

import (
	"context"
	"operator/output/fsm"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testcrdv1 "operator/api/v1"
)

const kind = "Testcrd"

var resource = &testcrdv1.Testcrd{
	TypeMeta: metav1.TypeMeta{
		Kind: kind,
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      resourceName.Name,
		Namespace: resourceName.Namespace,
	},
}

func setup(t *testing.T) {
	err := k8sClient.Create(context.TODO(), resource)
	require.NoError(t, err)
}

func teardown(t *testing.T) {
	err := k8sClient.Delete(context.TODO(), resource)
	require.NoError(t, err)

	resource = &testcrdv1.Testcrd{
		TypeMeta: metav1.TypeMeta{
			Kind: kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName.Name,
			Namespace: resourceName.Namespace,
		},
	}
}

// +vectorsigma:action:InitializeContext
func TestTestreconcileloop_InitializeContextAction(t *testing.T) {
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
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.InitializeContextAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Testreconcileloop.InitializeContextAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:LoadObjects
func TestTestreconcileloop_LoadObjectsAction(t *testing.T) {
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
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.LoadObjectsAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Testreconcileloop.LoadObjectsAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:SetReady
func TestTestreconcileloop_SetReadyAction(t *testing.T) {
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
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.SetReadyAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Testreconcileloop.SetReadyAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +vectorsigma:action:UpdateStatus
func TestTestreconcileloop_UpdateStatusAction(t *testing.T) {
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
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.UpdateStatusAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Testreconcileloop.UpdateStatusAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
