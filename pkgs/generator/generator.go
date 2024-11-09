/*
Copyright © 2024 Morten Hersson mhersson@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package generator

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/afero"

	"github.com/mhersson/vectorsigma/pkgs/shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
)

//go:embed templates/basic/*
var templates embed.FS

type InputParams struct {
	Input   string
	Output  string
	Module  string
	Package string
	Init    bool
}

type Generator struct {
	FS      afero.Fs
	Shell   shell.Interface
	FSM     *uml.FSM
	Module  string
	Package string
}

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// Get the current working directory name.
func GetCurrentWorkingDirectory(base bool) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	if base {
		return filepath.Base(dir), nil
	}

	return dir, nil
}

func (g *Generator) ExecuteTemplate(filename string) ([]byte, error) {
	tmpl, err := template.ParseFS(templates, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	buffer := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buffer, g)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.Bytes(), nil
}

// Initialize a new go module.
func (g *Generator) InitializeModule() error {
	err := g.Shell.NewCommand("go", "mod", "init", g.Module).Run()
	if err != nil {
		return fmt.Errorf("failed to initialize module: %w", err)
	}

	return nil
}

// Write file to disk.
func (g *Generator) WriteFile(path string, data []byte) error {
	err := afero.WriteFile(g.FS, path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Format the generated code.
func (g *Generator) FormatCode(path string) error {
	err := g.Shell.NewCommand("go", "fmt", path).Run()
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	return nil
}

func (g *Generator) Generate(init bool, input, output string) error {
	data, err := ReadFile(input)
	if err != nil {
		return err
	}

	g.FSM = uml.Parse(data)

	if init {
		if err := os.Chdir(output); err != nil {
			return fmt.Errorf("failed to change directory: %w", err)
		}

		if err := g.InitializeModule(); err != nil {
			return err
		}

		code, err := g.ExecuteTemplate("templates/basic/main.go.tmpl")
		if err != nil {
			return err
		}

		if err := g.WriteFile("main.go", code); err != nil {
			return err
		}

		if err := g.FS.MkdirAll(filepath.Join("pkg", g.Package), 0755); err != nil {
			return fmt.Errorf("failed to create package directory: %w", err)
		}

		output = filepath.Join("pkg", g.Package)
	} else {
		if err := g.FS.MkdirAll(filepath.Join(output, g.Package), 0755); err != nil {
			return fmt.Errorf("failed to create package directory: %w", err)
		}

		output = filepath.Join(output, g.Package)
	}

	for _, filename := range []string{"actions.go", "guards.go", "fsm.go", "state.go"} {
		code, err := g.ExecuteTemplate("templates/basic/" + filename + ".tmpl")
		if err != nil {
			return err
		}

		if filename == "fsm.go" {
			filename = "zz_generated_" + filename
		}

		if err := g.WriteFile(filepath.Join(output, filename), code); err != nil {
			return err
		}

		if err := g.FormatCode(filepath.Join(output, filename)); err != nil {
			return err
		}
	}

	return nil
}

func Run(input *InputParams) error {
	gen := &Generator{
		FS:    afero.NewOsFs(),
		Shell: &shell.Shell{},
	}

	var err error

	gen.Package = input.Package
	gen.Module = input.Module

	if input.Module == "" {
		gen.Module, err = GetCurrentWorkingDirectory(true)
		if err != nil {
			return err
		}
	}

	if input.Output == "" {
		input.Output, err = GetCurrentWorkingDirectory(false)
		if err != nil {
			return err
		}
	}

	return gen.Generate(input.Init, input.Input, input.Output)
}
