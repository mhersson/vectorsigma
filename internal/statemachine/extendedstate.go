package statemachine

import (
	"github.com/mhersson/vectorsigma/pkgs/generator"
)

type ExtendedState struct {
	Generator     *generator.Generator
	GeneratedData map[string][]byte
	Init          bool
	Input         string
	InputData     string
	Module        string
	Output        string
	Package       string
	Error         error
}
