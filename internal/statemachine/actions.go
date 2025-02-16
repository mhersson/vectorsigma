package statemachine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mhersson/vectorsigma/pkgs/generator"
	"github.com/mhersson/vectorsigma/pkgs/shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
	"github.com/spf13/afero"
)

// Actions

func (fsm *VectorSigma) InitializeAction(_ ...string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	if fsm.ExtendedState.Module == "" {
		fsm.ExtendedState.Module = filepath.Base(dir)
	}

	if fsm.ExtendedState.Output == "" {
		fsm.ExtendedState.Output = dir
	}

	fsm.Context.Generator = &generator.Generator{
		FS:      afero.NewOsFs(),
		Shell:   &shell.Shell{},
		Module:  fsm.ExtendedState.Module,
		Package: fsm.ExtendedState.Package,
	}

	fsm.ExtendedState.GeneratedData = make(map[string][]byte)

	return nil
}

func (fsm *VectorSigma) LoadInputAction(_ ...string) error {
	content, err := afero.ReadFile(fsm.Context.Generator.FS, fsm.ExtendedState.Input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	fsm.ExtendedState.InputData = string(content)

	return nil
}

func (fsm *VectorSigma) ExtractUMLAction(_ ...string) error {
	const startDelimiter = "```plantuml"

	const endDelimiter = "```"

	lenStartDelimiter := len(startDelimiter)
	markdown := fsm.ExtendedState.InputData

	startIndex := strings.Index(markdown, startDelimiter)
	if startIndex == -1 {
		return errors.New("no plantuml found in markdown")
	}

	endIndex := strings.Index(markdown[startIndex+lenStartDelimiter:], endDelimiter)
	if endIndex == -1 {
		return errors.New("missing end of plantuml code block in markdown")
	}

	fsm.ExtendedState.InputData = markdown[startIndex+lenStartDelimiter : endIndex+startIndex+lenStartDelimiter]

	return nil
}

func (fsm *VectorSigma) ParseUMLAction(_ ...string) error {
	fsm.Context.Generator.FSM = uml.Parse(fsm.ExtendedState.InputData)

	return nil
}

func (fsm *VectorSigma) GenerateStateMachineAction(_ ...string) error {
	files := []string{
		"actions.go",
		"actions_test.go",
		"guards.go",
		"guards_test.go",
		"statemachine.go",
		"extendedstate.go"}

	for _, filename := range files {
		code, err := fsm.Context.Generator.ExecuteTemplate("templates/application/" + filename + ".tmpl")
		if err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		if filename == "statemachine.go" {
			filename = "zz_generated_" + filename
		}

		fsm.ExtendedState.GeneratedData[filepath.Join(fsm.ExtendedState.Package, filename)] = code
	}

	return nil
}

func (fsm *VectorSigma) CreateOutputFolderAction(_ ...string) error {
	if err := fsm.Context.Generator.FS.MkdirAll(
		filepath.Join(fsm.ExtendedState.Output, fsm.ExtendedState.Package), 0755); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	return nil
}

func (fsm *VectorSigma) WriteGeneratedFilesAction(_ ...string) error {
	for filename, code := range fsm.ExtendedState.GeneratedData {
		if err := fsm.Context.Generator.WriteFile(filepath.Join(fsm.ExtendedState.Output, filename), code); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (fsm *VectorSigma) GenerateModuleFilesAction(_ ...string) error {
	files := []string{"main.go", "go.mod"}

	for _, filename := range files {
		code, err := fsm.Context.Generator.ExecuteTemplate("templates/application/" + filename + ".tmpl")
		if err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		fsm.ExtendedState.GeneratedData[filename] = code
	}

	return nil
}

func (fsm *VectorSigma) FormatCodeAction(_ ...string) error {
	for filename := range fsm.ExtendedState.GeneratedData {
		if filename == "go.mod" {
			continue
		}
		err := fsm.Context.Generator.FormatCode(filepath.Join(fsm.ExtendedState.Output, filename))
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}
