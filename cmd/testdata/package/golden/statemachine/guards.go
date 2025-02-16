package statemachine

func init() {
	AllGuards[IsError] = &IsErrorGuard{}
}

type IsErrorGuard struct{}

func (g *IsErrorGuard) Evaluate(state *ExtendedState) bool {
	//FIXME: Guard this anyway you want, but it should probably not look like this
	return false
}
