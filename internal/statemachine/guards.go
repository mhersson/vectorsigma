package statemachine

import "path/filepath"

// +vectorsigma:guard:IsError
func (fsm *VectorSigma) IsErrorGuard(_ ...string) bool {
	return fsm.ExtendedState.Error != nil
}

// +vectorsigma:guard:IsMarkdown
func (fsm *VectorSigma) IsMarkdownGuard(_ ...string) bool {
	return filepath.Ext(fsm.ExtendedState.Input) == ".md"
}

// +vectorsigma:guard:IsInitializingModule
func (fsm *VectorSigma) IsInitializingModuleGuard(_ ...string) bool {
	return fsm.ExtendedState.Init
}

// +vectorsigma:guard:PackageExists
func (fsm *VectorSigma) PackageExistsGuard(_ ...string) bool {
	return fsm.ExtendedState.PackageExists
}
