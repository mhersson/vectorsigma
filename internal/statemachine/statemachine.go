package statemachine

import (
	"log/slog"
	"os"
)

type (
	StateName  string
	ActionName string
	GuardName  string
)

const (
	FinalState                   StateName = "FinalState"
	Initializing                 StateName = "Initializing"
	LoadingInput                 StateName = "LoadingInput"
	ExtractingUML                StateName = "ExtractingUML"
	ParsingUML                   StateName = "ParsingUML"
	GeneratingStateMachine       StateName = "GeneratingStateMachine"
	CreatingInternalOutputFolder StateName = "CreatingInternalOutputFolder"
	CreatingOutputFolder         StateName = "CreatingOutputFolder"
	FilteringExistingFiles       StateName = "FilteringExistingFiles"
	MakingIncrementalUpdates     StateName = "MakingIncrementalUpdates"
	WritingGeneratedFiles        StateName = "WritingGeneratedFiles"
	GeneratingModuleFiles        StateName = "GeneratingModuleFiles"
	FormattingCode               StateName = "FormattingCode"
)

const (
	Initialize             ActionName = "Initialize"
	LoadInput              ActionName = "LoadInput"
	ExtractUML             ActionName = "ExtractUML"
	ParseUML               ActionName = "ParseUML"
	GenerateStateMachine   ActionName = "GenerateStateMachine"
	CreateOutputFolder     ActionName = "CreateOutPutFolder"
	FilterExistingFiles    ActionName = "FilterExistingFiles"
	MakeIncrementalUpdates ActionName = "MakeIncrementalUpdates"
	WriteGeneratedFiles    ActionName = "WriteGeneratedFiles"
	GenerateModuleFiles    ActionName = "GenerateModuleFiles"
	FormatCode             ActionName = "FormatCode"
)

const (
	IsError              GuardName = "IsError"
	IsMarkdown           GuardName = "IsMarkdown"
	IsInitializingModule GuardName = "IsInitializingModule"
	PackageExists        GuardName = "PackageExists"
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
		CurrentState:  Initializing,
		ExtendedState: &ExtendedState{},
		StateConfigs:  make(map[StateName]StateConfig),
	}

	// Define state configurations
	fsm.StateConfigs[Initializing] = StateConfig{
		Actions: []Action{
			{Name: Initialize, Execute: fsm.InitializeAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
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

	fsm.StateConfigs[ExtractingUML] = StateConfig{
		Actions: []Action{
			{Name: ExtractUML, Execute: fsm.ExtractUMLAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: ParsingUML,
		},
	}

	fsm.StateConfigs[ParsingUML] = StateConfig{
		Actions: []Action{
			{Name: ParseUML, Execute: fsm.ParseUMLAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: GeneratingStateMachine,
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

	fsm.StateConfigs[GeneratingModuleFiles] = StateConfig{
		Actions: []Action{
			{Name: GenerateModuleFiles, Execute: fsm.GenerateModuleFilesAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: CreatingInternalOutputFolder,
		},
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
			1: FilteringExistingFiles,
			2: WritingGeneratedFiles,
		},
	}

	fsm.StateConfigs[FilteringExistingFiles] = StateConfig{
		Actions: []Action{
			{Name: FilterExistingFiles, Execute: fsm.FilterExistingFilesAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: MakingIncrementalUpdates,
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
			1: WritingGeneratedFiles,
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

	fsm.StateConfigs[FormattingCode] = StateConfig{
		Actions: []Action{
			{Name: FormatCode, Execute: fsm.FormatCodeAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: FinalState,
		},
	}

	return fsm
}

// Run handles the state transitions based on the current state.
func (fsm *VectorSigma) Run() {
transitionsLoop:
	for {
		// If we are in the FinalState, exit the FSM
		if fsm.CurrentState == FinalState {
			return
		}

		config, exists := fsm.StateConfigs[fsm.CurrentState]

		if !exists {
			fsm.Context.Logger.Error("missing state config", "state", fsm.CurrentState)

			return
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
