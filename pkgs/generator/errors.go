package generator

type GenericError struct {
	Message string
}

func (g *GenericError) Error() string {
	return g.Message
}
