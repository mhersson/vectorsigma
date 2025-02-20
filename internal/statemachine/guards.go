package statemachine

import "path/filepath"

// Guards

// +vectorsigma:guard:IsError
func (fsm *VectorSigma) IsErrorGuard() bool {
	return fsm.ExtendedState.Error != nil
}

// +vectorsigma:guard:IsMarkdown
func (fsm *VectorSigma) IsMarkdownGuard() bool {
	return filepath.Ext(fsm.ExtendedState.Input) == ".md"
}

// +vectorsigma:guard:IsInitializingModule
func (fsm *VectorSigma) IsInitializingModuleGuard() bool {
	return fsm.ExtendedState.Init
}

// +vectorsigma:guard:PackageExists
func (fsm *VectorSigma) PackageExistsGuard() bool {
	return fsm.ExtendedState.PackageExits
}
