package statemachine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mhersson/vectorsigma/pkgs/generator"
	"github.com/mhersson/vectorsigma/pkgs/shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
	"github.com/spf13/afero"
)

// +vectorsigma:action:Initialize
func (fsm *VectorSigma) InitializeAction(_ ...string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	if fsm.ExtendedState.Module == "" {
		fsm.ExtendedState.Module = filepath.Base(dir)
	}

	relativePath := ""

	if fsm.ExtendedState.Output == "" {
		fsm.ExtendedState.Output = dir
	} else if fsm.ExtendedState.Output != "" && !fsm.ExtendedState.Init {
		// When generating a package for an existing project we need the relative path
		// to be able to correctly set the import in the actions and guards test files
		relativePath = fsm.ExtendedState.Output
		if strings.HasPrefix(relativePath, ".") || strings.HasPrefix(fsm.ExtendedState.Output, "/") {
			return errors.New("invalid output - output must be a sub directory of the current working directory without leading ./")
		}
	}

	if fsm.ExtendedState.Group == "" {
		fsm.ExtendedState.Group = fsm.ExtendedState.APIKind
	}

	fsm.Context.Generator = &generator.Generator{
		FS:           afero.NewOsFs(),
		Shell:        &shell.Shell{},
		APIKind:      fsm.ExtendedState.APIKind,
		APIVersion:   strings.ToLower(fsm.ExtendedState.APIVersion),
		Group:        strings.ToLower(fsm.ExtendedState.Group),
		Module:       fsm.ExtendedState.Module,
		Package:      fsm.ExtendedState.Package,
		Init:         fsm.ExtendedState.Init,
		RelativePath: relativePath,
		Version:      fsm.ExtendedState.VectorSigmaVersion,
	}

	fsm.ExtendedState.GeneratedFiles = make(map[string]GeneratedFile)

	return nil
}

// +vectorsigma:action:LoadInput
func (fsm *VectorSigma) LoadInputAction(_ ...string) error {
	content, err := afero.ReadFile(fsm.Context.Generator.FS, fsm.ExtendedState.Input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	fsm.ExtendedState.InputData = string(content)

	return nil
}

// +vectorsigma:action:ExtractUML
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

// +vectorsigma:action:ParseUML
func (fsm *VectorSigma) ParseUMLAction(_ ...string) error {
	fsm.Context.Generator.FSM = uml.Parse(fsm.ExtendedState.InputData)

	return nil
}

// +vectorsigma:action:GenerateStateMachine
func (fsm *VectorSigma) GenerateStateMachineAction(_ ...string) error {
	files := []string{
		"actions.go",
		"actions_test.go",
		"guards.go",
		"guards_test.go",
		"statemachine.go",
		"statemachine_test.go",
		"extendedstate.go",
	}

	templatePath := "templates/application"
	if fsm.ExtendedState.Operator {
		templatePath = "templates/operator"
	}

	for _, filename := range files {
		code, err := fsm.Context.Generator.ExecuteTemplate(filepath.Join(templatePath, filename+".tmpl"))
		if err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		if strings.HasPrefix(filename, "statemachine") {
			// Make it very clear that this is a generated file that should not be modified
			filename = "zz_generated_" + filename
		}

		fsm.ExtendedState.GeneratedFiles[filepath.Join(fsm.ExtendedState.Package, filename)] = GeneratedFile{Content: code, IncrementalChange: false}
	}

	return nil
}

// +vectorsigma:action:GenerateModuleFiles
func (fsm *VectorSigma) GenerateModuleFilesAction(_ ...string) error {
	files := []string{"main.go", "go.mod"}

	templatePath := "templates/application"
	if fsm.ExtendedState.Operator {
		templatePath = "templates/operator"
	}

	// Change the file path to projectroot/internal/package for new modules
	generatedFiles := make(map[string]GeneratedFile)
	for f, c := range fsm.ExtendedState.GeneratedFiles {
		generatedFiles[filepath.Join("internal", f)] = c
	}

	for _, filename := range files {
		if exists, err := fsm.Context.Generator.Exists(filepath.Join(fsm.ExtendedState.Output, filename)); exists || err != nil {
			if exists {
				return errors.New("failed to initialize new module. file exists " + filename)
			}

			if err != nil {
				return fmt.Errorf("failed to check if path exists %s - %w", filename, err)
			}
		}

		code, err := fsm.Context.Generator.ExecuteTemplate(filepath.Join(templatePath, filename+".tmpl"))
		if err != nil {
			return fmt.Errorf("code generation failed: %w", err)
		}

		generatedFiles[filename] = GeneratedFile{Content: code, IncrementalChange: false}
	}

	fsm.ExtendedState.GeneratedFiles = generatedFiles

	return nil
}

// +vectorsigma:action:CreateOutputFolder
func (fsm *VectorSigma) CreateOutputFolderAction(params ...string) error {
	outputfolder := filepath.Join(fsm.ExtendedState.Output, fsm.ExtendedState.Package)
	if len(params) > 0 {
		outputfolder = filepath.Join(fsm.ExtendedState.Output, params[0], fsm.ExtendedState.Package)
	}

	if exists, _ := fsm.Context.Generator.Exists(outputfolder); exists && !fsm.ExtendedState.Init {
		fsm.ExtendedState.PackageExits = true

		return nil
	}

	if err := fsm.Context.Generator.FS.MkdirAll(outputfolder, 0o755); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	return nil
}

// +vectorsigma:action:FilterGeneratedFiles
func (fsm *VectorSigma) FilterGeneratedFilesAction(_ ...string) error {
	// We should error out before if main.go or go.mod exists, and since package
	// exists so does probably extendtendstate.go too, but anyways..
	files := []string{"extendedstate.go", "main.go", "go.mod"}

	actionsAndguards := []string{"actions.go", "actions_test.go", "guards.go", "guards_test.go"}

	for filename, gf := range fsm.ExtendedState.GeneratedFiles {
		if exists, _ := fsm.Context.Generator.Exists(filepath.Join(fsm.ExtendedState.Output, filename)); exists {
			if slices.Contains(files, filepath.Base(filename)) {
				// extendedstate.go should never be overwritten
				delete(fsm.ExtendedState.GeneratedFiles, filename)
			}

			if slices.Contains(actionsAndguards, filepath.Base(filename)) && !gf.IncrementalChange {
				// don't write actions and guards unless they have changed
				delete(fsm.ExtendedState.GeneratedFiles, filename)
			}
		}
	}

	return nil
}

// +vectorsigma:action:MakeIncrementalUpdates
func (fsm *VectorSigma) MakeIncrementalUpdatesAction(_ ...string) error {
	files := []string{"actions.go", "actions_test.go", "guards.go", "guards_test.go"}

	for f, c := range fsm.ExtendedState.GeneratedFiles {
		if slices.Contains(files, filepath.Base(f)) {
			fullpath := filepath.Join(fsm.ExtendedState.Output, f)
			if exists, err := fsm.Context.Generator.Exists(fullpath); exists && err == nil {
				fsm.Context.Logger.Debug("Running incremental update", "file", f)

				code, changed, err := fsm.Context.Generator.IncrementalUpdate(fullpath, c.Content)
				if err != nil {
					return fmt.Errorf("incremental update failed: %w", err)
				}

				fsm.ExtendedState.GeneratedFiles[f] = GeneratedFile{Content: code, IncrementalChange: changed}
			} else if err != nil {
				return fmt.Errorf("failed to check if file exists: %w", err)
			}
		}
	}

	return nil
}

// +vectorsigma:action:WriteGeneratedFiles
func (fsm *VectorSigma) WriteGeneratedFilesAction(_ ...string) error {
	for filename, code := range fsm.ExtendedState.GeneratedFiles {
		if err := fsm.Context.Generator.WriteFile(filepath.Join(fsm.ExtendedState.Output, filename), code.Content); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// +vectorsigma:action:FormatCode
func (fsm *VectorSigma) FormatCodeAction(_ ...string) error {
	for filename := range fsm.ExtendedState.GeneratedFiles {
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
