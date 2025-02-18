package statemachine

import "path/filepath"

// Guards

func (fsm *VectorSigma) IsErrorGuard() bool {
	return fsm.ExtendedState.Error != nil
}

func (fsm *VectorSigma) IsMarkdownGuard() bool {
	return filepath.Ext(fsm.ExtendedState.Input) == ".md"
}

func (fsm *VectorSigma) IsStandaloneModuleGuard() bool {
	return fsm.ExtendedState.Init
}

func (fsm *VectorSigma) PackageExistsGuard() bool {
	return fsm.ExtendedState.PackageExits
}
