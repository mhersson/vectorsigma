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

type Generator struct {
	Module  string
	Package string
	FSM     *uml.FSM
	FS      afero.Fs
	Shell   shell.Interface
}

const ErrorPrefix = "\033[31mError:\033[0m"

func CheckError(err error) {
	if err != nil {
		fmt.Printf("%s %v\n", ErrorPrefix, err)

		os.Exit(1)
	}
}

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// Get the current working directory name.
func GetCurrentWorkingDirectory(base bool) string {
	dir, err := os.Getwd()
	CheckError(err)

	if base {
		return filepath.Base(dir)
	}

	return dir
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

func (g *Generator) Generate(init bool, input, output string) {
	data, err := ReadFile(input)
	CheckError(err)

	g.FSM = uml.Parse(data)

	if init {
		err := os.Chdir(output)
		CheckError(err)

		err = g.InitializeModule()
		CheckError(err)

		code, err := g.ExecuteTemplate("templates/basic/main.go.tmpl")
		CheckError(err)

		err = g.WriteFile("main.go", code)
		CheckError(err)

		err = g.FS.MkdirAll(filepath.Join("pkg", g.Package), 0755)
		CheckError(err)

		output = filepath.Join("pkg", g.Package)
	} else {
		err := g.FS.MkdirAll(filepath.Join(output, g.Package), 0755)
		CheckError(err)

		output = filepath.Join(output, g.Package)
	}

	for _, filename := range []string{"actions.go", "guards.go", "fsm.go", "state.go"} {
		code, err := g.ExecuteTemplate("templates/basic/" + filename + ".tmpl")
		CheckError(err)

		if filename == "fsm.go" {
			filename = "zz_generated_" + filename
		}

		err = g.WriteFile(filepath.Join(output, filename), code)
		CheckError(err)
		err = g.FormatCode(filepath.Join(output, filename))
		CheckError(err)
	}
}

func Run(init bool, module, pkg, input, output string) {
	gen := &Generator{
		FS:    afero.NewOsFs(),
		Shell: &shell.Shell{},
	}

	gen.Package = pkg
	gen.Module = module

	if module == "" {
		gen.Module = GetCurrentWorkingDirectory(true)
	}

	if output == "" {
		output = GetCurrentWorkingDirectory(false)
	}

	gen.Generate(init, input, output)
}
