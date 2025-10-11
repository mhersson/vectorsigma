/*
Copyright © 2024-2025 Morten Hersson <mhersson@gmail.com>

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
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/afero"

	"github.com/mhersson/vectorsigma/pkgs/shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
)

//go:embed templates/application/* templates/operator/*
var templates embed.FS

type Generator struct {
	FS           afero.Fs
	Shell        shell.Interface
	FSM          *uml.FSM
	APIKind      string
	APIVersion   string
	Group        string
	Module       string
	Package      string
	RelativePath string
	Version      string
	Init         bool
}

func (g *Generator) ExecuteTemplate(filename string) ([]byte, error) {
	titleTransformer := cases.Title(language.English)

	funcMap := template.FuncMap{
		"title":   titleTransformer.String,
		"toLower": strings.ToLower,
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

// Check if file or folder exists.
func (g *Generator) Exists(path string) (bool, error) {
	return afero.Exists(g.FS, path)
}

// Write file to disk.
func (g *Generator) WriteFile(path string, data []byte) error {
	err := afero.WriteFile(g.FS, path, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Format the generated code.
func (g *Generator) FormatCode(path string) error {
	const goCmd = "go"

	const goImportsCmd = "goimports"

	formatCmd := []string{goCmd, "fmt", path}
	if _, err := exec.LookPath(goImportsCmd); err == nil {
		formatCmd = []string{goImportsCmd, "-w", path}
	}

	if err := g.Shell.NewCommand(formatCmd[0], formatCmd[1:]...).Run(); err != nil {
		return fmt.Errorf("failed to format code at %s: %w", path, err)
	}

	return nil
}

// IncrementalUpdate compares the generated code with the existing code.  if any
// of the functions in the generated code are not in the existing code add them.
// If any functions are in the existing code and not in the generated code, and
// have a doc comment prefixed with '// +vectorsigma' remove them from the
// existing code. If any functions are in both trees, replace the existing code
// with the generated code if the existing code still contains a `// TODO:
// Impment me!` comment in the function body.
func (g *Generator) IncrementalUpdate(fullpath string, data []byte) ([]byte, bool, error) {
	// load existing code
	existing, err := afero.ReadFile(g.FS, fullpath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read existing code: %w", err)
	}

	containsChanges := false

	exisitingNode, err := decorator.Parse(existing)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse existing code: %w", err)
	}

	generatedNode, err := decorator.Parse(data)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse generated code: %w", err)
	}

	if changed := addOrReplace(exisitingNode, generatedNode); changed {
		containsChanges = true
	}

	if changed := removeNotInGenerated(exisitingNode, generatedNode); changed {
		containsChanges = true
	}

	var buf bytes.Buffer
	if err := decorator.Fprint(&buf, exisitingNode); err != nil {
		return nil, containsChanges, fmt.Errorf("failed to print modified code: %w", err)
	}

	return buf.Bytes(), containsChanges, nil
}

// addOrReplace compares the two files and if any of the functions in the
// generated code and not in the existing code add the new function to the
// existing code.  If any functions are in both trees, replace the existing code
// with the generated code if the existing code still contains a `// TODO:
// Impment me!` comment in the function body.
func addOrReplace(existingFile, generatedFile *dst.File) bool {
	containsChanges := false

	for _, genDecl := range generatedFile.Decls {
		if genDecl, ok := genDecl.(*dst.FuncDecl); ok {
			found := false

			for i, exDecl := range existingFile.Decls {
				exDecl, ok := exDecl.(*dst.FuncDecl)
				if !ok {
					continue
				}

				if exDecl.Name.Name == genDecl.Name.Name {
					found = true
					// Check if the existing function has a TODO comment
					hasTODO := false

					for _, line := range exDecl.Body.List {
						d := line.Decorations()
						if slices.Contains(d.Start.All(), "// TODO: Implement me!") {
							hasTODO = true

							break
						}
					}

					if hasTODO {
						// Replace the existing function with the generated one
						containsChanges = true
						existingFile.Decls[i] = genDecl
					}

					break
				}
			}

			if !found {
				// Add the new function to the existing code
				containsChanges = true

				existingFile.Decls = append(existingFile.Decls, genDecl)
			}
		}
	}

	return containsChanges
}

// removeNotInGenerated removes functions with the // +vectorsigma comment that
// are not in the generated code.
func removeNotInGenerated(exisitingNode, generatedNode *dst.File) bool {
	containsChanges := false

	for i := 0; i < len(exisitingNode.Decls); i++ {
		exDecl, ok := exisitingNode.Decls[i].(*dst.FuncDecl)
		if !ok {
			continue
		}

		found := false

		for _, genDecl := range generatedNode.Decls {
			genDecl, ok := genDecl.(*dst.FuncDecl)
			if !ok {
				continue
			}

			if exDecl.Name.Name == genDecl.Name.Name {
				found = true

				break
			}
		}

		if !found {
			for _, line := range exDecl.Decorations().Start {
				if strings.HasPrefix(line, "// +vectorsigma:") {
					containsChanges = true
					exisitingNode.Decls = slices.Delete(exisitingNode.Decls, i, i+1)
					i--

					break
				}
			}
		}
	}

	return containsChanges
}
