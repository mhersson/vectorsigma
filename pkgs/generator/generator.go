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
	Shell   *shell.Shell
}

const ErrorPrefix = "\033[31mError:\033[0m"

func CheckError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s %s, %v\n", ErrorPrefix, msg, err)
		os.Exit(1)
	}
}

func ReadFile(path string) string {
	content, err := os.ReadFile(path)
	CheckError(err, "failed to read file")

	return string(content)
}

// Get the current working directory name.
func GetCurrentWorkingDirectory(base bool) string {
	dir, err := os.Getwd()
	CheckError(err, "failed to get working directory")

	if base {
		return filepath.Base(dir)
	}

	return dir
}

func (g *Generator) ExecuteTemplate(filename string) []byte {
	tmpl, err := template.ParseFS(templates, "templates/"+filename)
	CheckError(err, "failed to parse template")

	buffer := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buffer, g)
	CheckError(err, "failed to execute template")

	return buffer.Bytes()
}

// Initialize a new go module.
func (g *Generator) InitializeModule() {
	err := g.Shell.NewCommand("go", "mod", "init", g.Module).Run()
	CheckError(err, "failed to initialize new go module")
}

// Write file to disk.
func (g *Generator) WriteFile(path string, data []byte) {
	err := afero.WriteFile(g.FS, path, data, 0644)
	CheckError(err, "failed to write file")
}

// Format the generated code.
func (g *Generator) FormatCode(path string) {
	err := g.Shell.NewCommand("go", "fmt", path).Run()
	CheckError(err, "failed to format code")
}

func (g *Generator) Generate(init bool, input, output string) {
	data := ReadFile(input)
	g.FSM = uml.Parse(data)

	if init {
		err := os.Chdir(output)
		CheckError(err, "failed to change directory")

		g.InitializeModule()

		code := g.ExecuteTemplate("basic/main.go.tmpl")
		g.WriteFile("main.go", code)

		err = g.FS.MkdirAll(filepath.Join("pkg", g.Package), 0755)
		CheckError(err, "failed create package directory")

		output = filepath.Join("pkg", g.Package)
	} else {
		err := g.FS.MkdirAll(filepath.Join(output, g.Package), 0755)
		CheckError(err, "failed create output directory")

		output = filepath.Join(output, g.Package)
	}

	for _, filename := range []string{"actions.go", "guards.go", "fsm.go", "state.go"} {
		code := g.ExecuteTemplate("basic/" + filename + ".tmpl")

		if filename == "fsm.go" {
			filename = "zz_generated_" + filename
		}

		g.WriteFile(filepath.Join(output, filename), code)
		g.FormatCode(filepath.Join(output, filename))
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
