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
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/mhersson/vectorsigma/pkgs/shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
)

//go:embed templates/application/*
var templates embed.FS

type Generator struct {
	FS      afero.Fs
	Shell   shell.Interface
	FSM     *uml.FSM
	Module  string
	Package string
}

func (g *Generator) ExecuteTemplate(filename string) ([]byte, error) {
	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
	}

	tmpl, err := template.New(filepath.Base(filename)).Funcs(funcMap).ParseFS(templates, filename)
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
