package statemachine

import ()

func init() {
	AllActions[SwitchIn] = &SwitchInAction{}
}

type SwitchInAction struct{}

func (a *SwitchInAction) Execute(ctx *Context, state *ExtendedState, parameters ...string) error {
	// FIXME: Fake it or break it,  but don't leave it like this
	return nil
}
