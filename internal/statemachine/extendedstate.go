package statemachine

import (
	"log/slog"

	"github.com/mhersson/vectorsigma/pkgs/generator"
)

// A struct that holds the items needed for the actions to do their work.
// Things like client libraries and loggers, go here.
type Context struct {
	Logger    *slog.Logger // Do NOT delete this!
	Generator *generator.Generator
}

// A struct that holds the "extended state" of the state machine, including data
// being fetched and read. This should only be modified by actions, while guards
// should only read the extended state to assess their value.
type ExtendedState struct {
	GeneratedFiles map[string]GeneratedFile
	Init           bool
	PackageExits   bool
	Operator       bool
	APIVersion     string
	APIKind        string
	Input          string
	InputData      string
	Module         string
	Output         string
	Package        string
	Error          error
}

type GeneratedFile struct {
	Content           []byte
	IncrementalChange bool
}
