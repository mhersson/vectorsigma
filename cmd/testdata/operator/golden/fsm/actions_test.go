package fsm_test

import (
	"context"
	"operator/output/fsm"
	"testing"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	unitv1 "operator/api/v1"
)

const kind = "TestCRD"

// resourceName is used by both unit and integration tests
var resourceName = types.NamespacedName{
	Namespace: "default",
	Name:      "test-resource",
}

// fakeK8sClient is a shared variable that can be used by tests
// Unit tests create their own fake clients per test
// Integration tests will use the envtest k8sClient
var fakeK8sClient client.Client

// silentLogger creates a logger that discards all output
func silentLogger() logr.Logger {
	return logr.Discard()
}

// testContext returns a fully configured test context
func testContext() *fsm.Context {
	// Create a new scheme and register all types we might need in tests
	testScheme := scheme.Scheme
	_ = unitv1.AddToScheme(testScheme)

	// Create a fake client with the comprehensive scheme and status subresource support
	fakeClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithStatusSubresource(&unitv1.TestCRD{}).
		Build()

	return &fsm.Context{
		Logger: silentLogger(),
		Client: fakeClient,
		Ctx:    context.TODO(),
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
