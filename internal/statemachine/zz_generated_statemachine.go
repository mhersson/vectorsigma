// This file is generated by VectorSigma. DO NOT EDIT.
package statemachine

import (
	"fmt"
	"log/slog"
	"os"
)

type (
	StateName  string
	ActionName string
	GuardName  string
)

const (
	CreatingInternalOutputFolder StateName = "CreatingInternalOutputFolder"
	CreatingOutputFolder         StateName = "CreatingOutputFolder"
	ExtractingUML                StateName = "ExtractingUML"
	FilteringGeneratedFiles      StateName = "FilteringGeneratedFiles"
	FinalState                   StateName = "FinalState"
	FormattingCode               StateName = "FormattingCode"
	GeneratingModuleFiles        StateName = "GeneratingModuleFiles"
	GeneratingStateMachine       StateName = "GeneratingStateMachine"
	InitialState                 StateName = "InitialState"
	Initializing                 StateName = "Initializing"
	LoadingInput                 StateName = "LoadingInput"
	MakingIncrementalUpdates     StateName = "MakingIncrementalUpdates"
	ParsingUML                   StateName = "ParsingUML"
	WritingGeneratedFiles        StateName = "WritingGeneratedFiles"
)

const (
	CreateOutputFolder     ActionName = "CreateOutputFolder"
	ExtractUML             ActionName = "ExtractUML"
	FilterGeneratedFiles   ActionName = "FilterGeneratedFiles"
	FormatCode             ActionName = "FormatCode"
	GenerateModuleFiles    ActionName = "GenerateModuleFiles"
	GenerateStateMachine   ActionName = "GenerateStateMachine"
	Initialize             ActionName = "Initialize"
	LoadInput              ActionName = "LoadInput"
	MakeIncrementalUpdates ActionName = "MakeIncrementalUpdates"
	ParseUML               ActionName = "ParseUML"
	WriteGeneratedFiles    ActionName = "WriteGeneratedFiles"
)

const (
	IsError              GuardName = "IsError"
	IsInitializingModule GuardName = "IsInitializingModule"
	IsMarkdown           GuardName = "IsMarkdown"
	PackageExists        GuardName = "PackageExists"
)

const maxStateDepth = 5

// Action represents a function that can be executed in a state and may return an error.
type Action struct {
	Name    ActionName
	Params  []string
	Execute func(...string) error
}

// Guard represents a function that returns a boolean indicating if a transition should occur.
type Guard struct {
	Name   GuardName
	Check  func() bool
	Action *Action
}

// StateConfig holds the actions and guards for a state.
type StateConfig struct {
	Actions     []Action
	Guards      []Guard
	Transitions map[int]StateName // Maps guard index to the next state
	Composite   CompositeState
}

type CompositeState struct {
	InitialState StateName
	StateConfigs map[StateName]StateConfig
}

// VectorSigma represents the Finite State Machine (fsm) for VectorSigma.
type VectorSigma struct {
	Context       *Context
	CurrentState  StateName
	ExtendedState *ExtendedState
	StateConfigs  map[StateName]StateConfig
}

// New initializes a new FSM.
func New() *VectorSigma {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)

	if os.Getenv("VECTORSIGMA_DEBUG") != "" {
		logLevel.Set(slog.LevelDebug)
	}

	fsm := &VectorSigma{
		Context:       &Context{Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))},
		CurrentState:  InitialState,
		ExtendedState: &ExtendedState{},
		StateConfigs:  make(map[StateName]StateConfig),
	}
	fsm.StateConfigs[CreatingInternalOutputFolder] = StateConfig{
		Actions: []Action{
			{Name: CreateOutputFolder, Execute: fsm.CreateOutputFolderAction, Params: []string{"internal"}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: WritingGeneratedFiles,
		},
	}
	fsm.StateConfigs[CreatingOutputFolder] = StateConfig{
		Actions: []Action{
			{Name: CreateOutputFolder, Execute: fsm.CreateOutputFolderAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
			{Name: PackageExists, Check: fsm.PackageExistsGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: MakingIncrementalUpdates,
			2: WritingGeneratedFiles,
		},
	}
	fsm.StateConfigs[ExtractingUML] = StateConfig{
		Actions: []Action{
			{Name: ExtractUML, Execute: fsm.ExtractUMLAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: ParsingUML,
		},
	}
	fsm.StateConfigs[FilteringGeneratedFiles] = StateConfig{
		Actions: []Action{
			{Name: FilterGeneratedFiles, Execute: fsm.FilterGeneratedFilesAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: WritingGeneratedFiles,
		},
	}

	fsm.StateConfigs[FormattingCode] = StateConfig{
		Actions: []Action{
			{Name: FormatCode, Execute: fsm.FormatCodeAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: FinalState,
		},
	}
	fsm.StateConfigs[GeneratingModuleFiles] = StateConfig{
		Actions: []Action{
			{Name: GenerateModuleFiles, Execute: fsm.GenerateModuleFilesAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: CreatingInternalOutputFolder,
		},
	}
	fsm.StateConfigs[GeneratingStateMachine] = StateConfig{
		Actions: []Action{
			{Name: GenerateStateMachine, Execute: fsm.GenerateStateMachineAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
			{Name: IsInitializingModule, Check: fsm.IsInitializingModuleGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: GeneratingModuleFiles,
			2: CreatingOutputFolder,
		},
	}
	fsm.StateConfigs[InitialState] = StateConfig{
		Actions: []Action{},
		Guards:  []Guard{},
		Transitions: map[int]StateName{
			0: Initializing,
		},
	}
	fsm.StateConfigs[Initializing] = StateConfig{
		Actions: []Action{
			{Name: Initialize, Execute: fsm.InitializeAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: LoadingInput,
		},
	}
	fsm.StateConfigs[LoadingInput] = StateConfig{
		Actions: []Action{
			{Name: LoadInput, Execute: fsm.LoadInputAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
			{Name: IsMarkdown, Check: fsm.IsMarkdownGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: ExtractingUML,
			2: ParsingUML,
		},
	}
	fsm.StateConfigs[MakingIncrementalUpdates] = StateConfig{
		Actions: []Action{
			{Name: MakeIncrementalUpdates, Execute: fsm.MakeIncrementalUpdatesAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: FilteringGeneratedFiles,
		},
	}
	fsm.StateConfigs[ParsingUML] = StateConfig{
		Actions: []Action{
			{Name: ParseUML, Execute: fsm.ParseUMLAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: GeneratingStateMachine,
		},
	}
	fsm.StateConfigs[WritingGeneratedFiles] = StateConfig{
		Actions: []Action{
			{Name: WriteGeneratedFiles, Execute: fsm.WriteGeneratedFilesAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: FormattingCode,
		},
	}

	return fsm
}

// Run handles the state transitions based on the current state.
func (fsm *VectorSigma) Run() error {
	return run(fsm, fsm.StateConfigs, 0)
}

func run(fsm *VectorSigma, stateConfigs map[StateName]StateConfig, depth int) error {
	if depth > maxStateDepth {
		return fmt.Errorf("max state depth exceeded")
	}

	for {
		// If we are in the FinalState, exit the FSM
		if fsm.CurrentState == FinalState {
			// Reset to the Initial State in case the FSM is run in a loop
			fsm.CurrentState = InitialState

			return fsm.ExtendedState.Error
		}

		config, exists := stateConfigs[fsm.CurrentState]

		if !exists {
			fsm.Context.Logger.Error("missing config", "state", fsm.CurrentState)

			return fmt.Errorf("missing config for state: %s", fsm.CurrentState)
		}

		if config.Composite.StateConfigs != nil {
			parentState := fsm.CurrentState
			// Recursively run the composite state machine
			fsm.CurrentState = config.Composite.InitialState
			fsm.Context.Logger.Debug("entering composite state", "state", parentState, "initial", fsm.CurrentState)
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
		nextState, err := runAllGuards(fsm.Context, fsm.CurrentState, config)
		if err != nil {
			// Guarded actions will always transition to the FinalState
			fsm.ExtendedState.Error = err
			fsm.CurrentState = FinalState
		}
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

func runAllGuards(context *Context, currentState StateName, config StateConfig) (StateName, error) {
	for guardIndex, guard := range config.Guards {
		if guard.Check() {
			if guard.Action != nil {
				action := guard.Action
				if err := action.Execute(action.Params...); err != nil {
					context.Logger.Debug("guarded action failed", "state", currentState,
						"guard", guard.Name, "action", action.Name, "error", err)

					return "", err
				}
			}

			// Transition to the state mapped to this guard index
			if nextState, exists := config.Transitions[guardIndex]; exists {
				context.Logger.Debug("guarded transition", "guard", guard.Name, "current", currentState, "next", nextState)

				return nextState, nil
			}
		}
	}

	return "", nil
}
