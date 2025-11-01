package fsm

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	unitv1 "operator/api/v1"
)

// A struct that holds the items needed for the actions to do their work.
// Things like client libraries and loggers, go here.
type Context struct {
	Logger   logr.Logger
	Client   client.Client
	Ctx      context.Context
	Recorder record.EventRecorder
}

// A struct that holds the "extended state" of the state machine, including data
// being fetched and read. This should only be modified by actions, while guards
// should only read the extended state to assess their value.
type ExtendedState struct {
	Error        error
	Result       ctrl.Result
	ResourceName types.NamespacedName
	Instance     unitv1.TestCRD
}
