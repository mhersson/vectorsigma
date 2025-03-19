// This file is generated by VectorSigma. DO NOT EDIT.
package fsm

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
)

type (
	StateName  string
	ActionName string
	GuardName  string
)

const (
	FinalState          StateName = "FinalState"
	InitialState        StateName = "InitialState"
	InitializingContext StateName = "InitializingContext"
	LoadingObjects      StateName = "LoadingObjects"
	SettingReady        StateName = "SettingReady"
	UpdatingStatus      StateName = "UpdatingStatus"
)

const (
	InitializeContext ActionName = "InitializeContext"
	LoadObjects       ActionName = "LoadObjects"
	SetReady          ActionName = "SetReady"
	UpdateStatus      ActionName = "UpdateStatus"
)

const (
	IsError  GuardName = "IsError"
	NotFound GuardName = "NotFound"
)

// Action represents a function that can be executed in a state and may return an error.
type Action struct {
	Name    ActionName
	Params  []string
	Execute func(...string) error
}

// Guard represents a function that returns a boolean indicating if a transition should occur.
type Guard struct {
	Name  GuardName
	Check func() bool
}

// StateConfig holds the actions and guards for a state.
type StateConfig struct {
	Actions     []Action
	Guards      []Guard
	Transitions map[int]StateName // Maps guard index to the next state
}

// VectorSigma represents the Finite State Machine (fsm) for VectorSigma.
type Testreconcileloop struct {
	Context       *Context
	CurrentState  StateName
	ExtendedState *ExtendedState
	StateConfigs  map[StateName]StateConfig
}

// New initializes a new FSM.
func New() *Testreconcileloop {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	if os.Getenv("TESTRECONCILELOOP_DEBUG") != "" {
		logLevel.Set(slog.LevelDebug)
	}

	fsm := &Testreconcileloop{
		Context:       &Context{Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))},
		CurrentState:  InitialState,
		ExtendedState: &ExtendedState{},
		StateConfigs:  make(map[StateName]StateConfig),
	}

	fsm.StateConfigs[InitialState] = StateConfig{
		Actions: []Action{},
		Guards:  []Guard{},
		Transitions: map[int]StateName{
			0: InitializingContext,
		},
	}
	fsm.StateConfigs[InitializingContext] = StateConfig{
		Actions: []Action{
			{Name: InitializeContext, Execute: fsm.InitializeContextAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: LoadingObjects,
		},
	}
	fsm.StateConfigs[LoadingObjects] = StateConfig{
		Actions: []Action{
			{Name: LoadObjects, Execute: fsm.LoadObjectsAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
			{Name: NotFound, Check: fsm.NotFoundGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: FinalState,
			2: SettingReady,
		},
	}
	fsm.StateConfigs[SettingReady] = StateConfig{
		Actions: []Action{
			{Name: SetReady, Execute: fsm.SetReadyAction, Params: []string{}},
		},
		Guards: []Guard{},
		Transitions: map[int]StateName{
			0: UpdatingStatus,
		},
	}
	fsm.StateConfigs[UpdatingStatus] = StateConfig{
		Actions: []Action{
			{Name: UpdateStatus, Execute: fsm.UpdateStatusAction, Params: []string{}},
		},
		Guards: []Guard{},
		Transitions: map[int]StateName{
			0: FinalState,
		},
	}

	return fsm
}

// Run handles the state transitions based on the current state.
func (fsm *Testreconcileloop) Run() (ctrl.Result, error) {
transitionsLoop:
	for {
		// If we are in the FinalState, exit the FSM
		if fsm.CurrentState == FinalState {
			return fsm.ExtendedState.Result, fsm.ExtendedState.Error
		}

		config, exists := fsm.StateConfigs[fsm.CurrentState]

		if !exists {
			fsm.Context.Logger.Error("missing state config", "state", fsm.CurrentState)

			return ctrl.Result{}, errors.New("missing state config for " + string(fsm.CurrentState))
		}

		// Execute all actions for the current state
		for _, action := range config.Actions {
			fsm.Context.Logger.Debug("executing", "action", action.Name, "state", fsm.CurrentState)

			if err := action.Execute(action.Params...); err != nil {
				fsm.Context.Logger.Error("action failed", "action", action.Name, "state", fsm.CurrentState, "error", err)
				fsm.ExtendedState.Error = err

				break
			}
		}

		// Check guards and determine the next state
		for guardIndex, guard := range config.Guards {
			if guard.Check() {
				// Transition to the state mapped to this guard index
				if nextState, exists := config.Transitions[guardIndex]; exists {
					fsm.Context.Logger.Debug("guarded transition", "guard", guard.Name, "current", fsm.CurrentState, "next", nextState)

					fsm.CurrentState = nextState

					continue transitionsLoop
				}
			}
		}
		// Check for unguarded transition
		if nextState, exists := config.Transitions[len(config.Guards)]; exists {
			fsm.Context.Logger.Debug("unguarded transition", "current", fsm.CurrentState, "next", nextState)
			fsm.CurrentState = nextState
		}
	}
}

func (fsm *Testreconcileloop) Run() (ctrl.Result, error) {
	return run(fsm, fsm.StateConfigs, 0)
}

func run(fsm *Testreconcileloop, stateConfigs map[StateName]StateConfig, depth int) (ctrl.Result, error) {
	if depth > MaxStateDepth {
		return ctrl.Result{}, fmt.Errorf("max state depth exceeded")
	}

	for {
		// If we are in the FinalState, exit the FSM
		if fsm.CurrentState == FinalState {
			// Reset to the Initial State in case the FSM is run in a loop
			fsm.CurrentState = InitialState

			return fsm.ExtendedState.Result, fsm.ExtendedState.Error
		}

		config, exists := stateConfigs[fsm.CurrentState]

		if !exists {
			fsm.Context.Logger.Error("missing config", "state", fsm.CurrentState)

			return ctrl.Result{}, errors.New("missing state config for " + string(fsm.CurrentState))
		}

		if config.Composite.StateConfigs != nil {
			parentState := fsm.CurrentState
			// Recursively run the compound state machine
			fsm.CurrentState = config.Composite.InitialState
			fsm.Context.Logger.Debug("entering compound state", "state", parentState, "initial", fsm.CurrentState)
			err := run(fsm, config.Composite.StateConfigs, depth+1)
			if err != nil {
				fsm.Context.Logger.Error("composite state machine failed", "state", fsm.CurrentState, "error", err)
				fsm.ExtendedState.Error = err
			}

			fsm.Context.Logger.Debug("exiting composite state", "state", parentState)
			fsm.CurrentState = parentState
		} else {
			// Execute all actions for the current state
			err := runAllActions(fsm.Context, fsm.CurrentState, config.Actions)
			if err != nil {
				fsm.Context.Logger.Error("action failed", "state", fsm.CurrentState, "error", err)
				fsm.ExtendedState.Error = err
			}
		}

		// Check guards and determine the next state
		nextState := runAllGuards(fsm.Context, fsm.CurrentState, config)
		if nextState != "" {
			fsm.CurrentState = nextState

			continue
		}

		// Check for unguarded transition
		if nextState, exists := config.Transitions[len(config.Guards)]; exists {
			fsm.Context.Logger.Debug("unguarded transition", "current", fsm.CurrentState, "next", nextState)
			fsm.CurrentState = nextState
		}
	}

}

func runAllActions(context *Context, currentState StateName, actions []Action) error {
	for _, action := range actions {
		context.Logger.Debug("executing", "action", action.Name, "state", currentState)

		if err := action.Execute(action.Params...); err != nil {
			return err
		}
	}

	return nil
}

func runAllGuards(context *Context, currentState StateName, config StateConfig) StateName {
	for guardIndex, guard := range config.Guards {
		if guard.Check() {
			// Transition to the state mapped to this guard index
			if nextState, exists := config.Transitions[guardIndex]; exists {
				context.Logger.Debug("guarded transition", "guard", guard.Name, "current", currentState, "next", nextState)

				return nextState
			}
		}
	}

	return ""
}
