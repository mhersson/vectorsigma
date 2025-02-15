package statemachine

// Guards

func (fsm *FSM) IsErrorGuard() bool {
	return false
}

func (fsm *FSM) IsMarkdownGuard() bool {
	return false
}

func (fsm *FSM) IsStandaloneModuleGuard() bool {
	return false
}
