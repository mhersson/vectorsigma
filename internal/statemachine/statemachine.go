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
	FinalState             StateName = "FinalState"
	Initializing           StateName = "Initializing"
	LoadingInput           StateName = "LoadingInput"
	ExtractingUML          StateName = "ExtractingUML"
	ParsingUML             StateName = "ParsingUML"
	GeneratingStateMachine StateName = "GeneratingStateMachine"
	CreatingOutputFolder   StateName = "CreatingOutputFolder"
	WritingGeneratedFiles  StateName = "WritingGeneratedFiles"
	InitializingGoModule   StateName = "InitializingGoModule"
	GeneratingMainFile     StateName = "GeneratingMainFile"
	FormattingCode         StateName = "FormattingCode"
)

const (
	Initialize           ActionName = "Initialize"
	LoadInput            ActionName = "LoadInput"
	ExtractUML           ActionName = "ExtractUML"
	ParseUML             ActionName = "ParseUML"
	GenerateStateMachine ActionName = "GenerateStateMachine"
	CreateOutputFolder   ActionName = "CreateOutPutFolder"
	WriteGeneratedFiles  ActionName = "WriteGeneratedFiles"
	InitializeGoModule   ActionName = "InitializingGoModule"
	GenerateMainFile     ActionName = "GenerateMainFile"
	FormatCode           ActionName = "FormatCode"
)

const (
	IsError            GuardName = "IsError"
	IsMarkdown         GuardName = "IsMarkdown"
	IsStandaloneModule GuardName = "IsStandaloneModule"
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

// FSM represents the FSM for VectorSigma.
type FSM struct {
	Logger        *slog.Logger
	CurrentState  StateName
	StateConfigs  map[StateName]StateConfig
	ExtendedState *ExtendedState
}

// New initializes a new finite state machine.
func New() *FSM {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	fsm := &FSM{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{

			Level: logLevel,
		})),
		CurrentState:  Initializing,
		StateConfigs:  make(map[StateName]StateConfig),
		ExtendedState: &ExtendedState{},
	}

	// Define state configurations
	fsm.StateConfigs[Initializing] = StateConfig{
		Actions: []Action{
			{Name: Initialize, Execute: fsm.InitializeAction, Params: []string{}},
		},
		Guards: []Guard{},
		Transitions: map[int]StateName{
			0: LoadingInput,
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
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: CreatingOutputFolder,
		},
	}

	fsm.StateConfigs[CreatingOutputFolder] = StateConfig{
		Actions: []Action{
			{Name: CreateOutputFolder, Execute: fsm.CreateOutputFolderAction, Params: []string{}},
		},
		Guards: []Guard{
			{Name: IsError, Check: fsm.IsErrorGuard},
			{Name: IsStandaloneModule, Check: fsm.IsStandaloneModuleGuard},
		},
		Transitions: map[int]StateName{
			0: FinalState,
			1: InitializingGoModule,
			2: WritingGeneratedFiles,
		},
	}

	fsm.StateConfigs[InitializingGoModule] = StateConfig{
		Actions: []Action{
			{Name: InitializeGoModule, Execute: fsm.InitializeGoModuleAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
		Transitions: map[int]StateName{
			0: FinalState,
			1: GeneratingMainFile,
		},
	}

	fsm.StateConfigs[GeneratingMainFile] = StateConfig{
		Actions: []Action{
			{Name: GenerateMainFile, Execute: fsm.GenerateMainFileAction, Params: []string{}},
		},
		Guards: []Guard{{Name: IsError, Check: fsm.IsErrorGuard}},
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
func (fsm *FSM) Run() {
transitionsLoop:
	for {
		// If we are in the FinalState, exit the FSM
		if fsm.CurrentState == FinalState {
			return
		}

		config, exists := fsm.StateConfigs[fsm.CurrentState]

		if !exists {
			fsm.Logger.Error("missing state config", "state", fsm.CurrentState)

			return
		}

		// Execute all actions for the current state
		for _, action := range config.Actions {
			fsm.Logger.Debug("executing", "action", action.Name, "state", fsm.CurrentState)

			if err := action.Execute(action.Params...); err != nil {
				fsm.Logger.Error("action failed", "action", action.Name, "state", fsm.CurrentState, "error", err)
				// FIXME: This migt not always be the wanted outcome
				fsm.CurrentState = FinalState

				break
			}
		}

		// Check guards and determine the next state
		if fsm.CurrentState != FinalState {
			for guardIndex, guard := range config.Guards {
				if guard.Check() {
					// Transition to the state mapped to this guard index
					if nextState, exists := config.Transitions[guardIndex]; exists {
						fsm.Logger.Debug("guarded transition", "guard", guard.Name, "current", fsm.CurrentState, "next", nextState)

						fsm.CurrentState = nextState

						continue transitionsLoop
					}
				}
			}
			// Check for unguarded transition
			if nextState, exists := config.Transitions[len(config.Guards)]; exists {
				fsm.Logger.Debug("unguarded transition", "current", fsm.CurrentState, "next", nextState)
				fsm.CurrentState = nextState
			}
		}
	}
}
