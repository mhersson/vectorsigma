package statemachine

import "path/filepath"

// Guards

func (fsm *FSM) IsErrorGuard() bool {
	return false
}

func (fsm *FSM) IsMarkdownGuard() bool {
	return filepath.Ext(fsm.ExtendedState.Input) == ".md"
}

func (fsm *FSM) IsStandaloneModuleGuard() bool {
	return fsm.ExtendedState.Init
}
