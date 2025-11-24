package fsm

// +vectorsigma:guard:IsError
func (fsm *Testreconcileloop) IsErrorGuard(_ ...string) bool {
	// TODO: Implement me!
	return false
}

// +vectorsigma:guard:NotFound
func (fsm *Testreconcileloop) NotFoundGuard(_ ...string) bool {
	// TODO: Implement me!
	return false
}
