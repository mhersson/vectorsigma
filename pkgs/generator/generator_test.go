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
package generator_test

import (
	"os/exec"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mhersson/vectorsigma/pkgs/generator"
	mock_shell "github.com/mhersson/vectorsigma/pkgs/shell/mock_shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
)

func TestGenerator_ExecuteTemplate(t *testing.T) {
	tests := []struct {
		name      string
		generator *generator.Generator
		filename  string
		want      string
		wantErr   bool
	}{
		{
			name: "ValidTemplate",
			generator: &generator.Generator{
				Package: "statemachine",
				FS:      afero.NewMemMapFs(),
				FSM: &uml.FSM{
					GuardNames: []string{"IsError"},
					Title:      "UnitTest",
				},
			},
			filename: "templates/application/guards.go.tmpl",
			wantErr:  false,
			want: `package statemachine
// +vectorsigma:guard:IsError
func (fsm *UnitTest) IsErrorGuard(_ ...string) bool {
	// TODO: Implement me!
	return false
}
`,
		},
		{
			name: "Template does not exist",
			generator: &generator.Generator{
				Package: "statemachine",
				FS:      afero.NewMemMapFs(),
				FSM: &uml.FSM{
					GuardNames: []string{"IsError"},
					Title:      "UnitTest",
				},
			},
			filename: "templates/application/does_not_exist.go.tmpl",
			wantErr:  true,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, err := tt.generator.ExecuteTemplate(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("ExecuteTemplate = %v, want %v", got, tt.want)
			} else if !tt.wantErr {
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestGenerator_Exists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *generator.Generator
		filepath string
		exists   bool
	}{
		{
			name: "FileDoesNotExist",
			setup: func() *generator.Generator {
				return &generator.Generator{FS: afero.NewMemMapFs()}
			},
			filepath: "/path/to/file",
			exists:   false,
		},
		{
			name: "FileExists",
			setup: func() *generator.Generator {
				fs := afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "/path/to/file", []byte("content"), 0o644)

				return &generator.Generator{FS: fs}
			},
			filepath: "/path/to/file",
			exists:   true,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := tt.setup()
			exists, err := g.Exists(tt.filepath)
			require.NoError(t, err)
			assert.Equal(t, tt.exists, exists)
		})
	}
}

func TestGenerator_WriteFile(t *testing.T) {
	tests := []struct {
		name      string
		generator *generator.Generator
		filepath  string
		content   []byte
	}{
		{
			name:      "WriteFile",
			generator: &generator.Generator{FS: afero.NewMemMapFs()},
			filepath:  "/path/to/file",
			content:   []byte("content"),
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.generator.WriteFile(tt.filepath, tt.content)
			require.NoError(t, err)

			content, err := afero.ReadFile(tt.generator.FS, tt.filepath)
			require.NoError(t, err)
			assert.Equal(t, tt.content, content)
		})
	}
}

func TestGenerator_FormatCode(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *generator.Generator
		filepath string
	}{
		{
			name: "FormatCode",
			setup: func() *generator.Generator {
				mockShell := mock_shell.NewMockInterface(t)
				mockCmd := mock_shell.NewMockCmdRunner(t)

				g := &generator.Generator{
					FS:    afero.NewMemMapFs(),
					Shell: mockShell,
				}

				if _, err := exec.LookPath("goimports"); err == nil {
					mockShell.On("NewCommand", "goimports", "-w", "/path/to/file").Return(mockCmd)
					mockCmd.EXPECT().Run().Return(nil)
				} else {
					mockShell.On("NewCommand", "go", "fmt", "/path/to/file").Return(mockCmd)
					mockCmd.EXPECT().Run().Return(nil)
				}

				return g
			},
			filepath: "/path/to/file",
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := tt.setup()
			err := g.FormatCode(tt.filepath)
			require.NoError(t, err)
			g.Shell.(*mock_shell.MockInterface).AssertExpectations(t)
		})
	}
}

func TestGenerator_IncrementalUpdate(t *testing.T) {
	tests := []struct {
		name          string
		generator     *generator.Generator
		filepath      string
		existingCode  string
		generatedCode string
		wantedCode    string
		changed       bool
	}{
		{
			name:      "IncrementalUpdate",
			generator: &generator.Generator{FS: afero.NewMemMapFs()},
			filepath:  "/path/to/file.go",
			existingCode: `package statemachine

// +vectorsigma:action:SwitchIn
func (fsm *TrafficLight) SwitchInAction(_ ...string) (string, error) {
	// TODO: Implement me!
	// This should be overwritten
	return "", nil
}

// +vectorsigma:action:helloworld
func helloworld() error {
	// this should be deleted
	return nil
}

func thisshouldbeleftalone() error {
	return nil
}
`,

			generatedCode: `package statemachine

// +vectorsigma:action:SwitchIn
func (fsm *TrafficLight) SwitchInAction(_ ...string) error {
	// TODO: Implement me!
	return nil
}

// +vectorsigma:action:THIS_SHOULD_BE_ADDED
func (fsm *TrafficLight) AddMe(_ ...string) error {
	// TODO: Implement me!
	return nil
}
`,

			wantedCode: `package statemachine

// +vectorsigma:action:SwitchIn
func (fsm *TrafficLight) SwitchInAction(_ ...string) error {
	// TODO: Implement me!
	return nil
}

func thisshouldbeleftalone() error {
	return nil
}

// +vectorsigma:action:THIS_SHOULD_BE_ADDED
func (fsm *TrafficLight) AddMe(_ ...string) error {
	// TODO: Implement me!
	return nil
}
`,
			changed: true,
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_ = afero.WriteFile(tt.generator.FS, tt.filepath, []byte(tt.existingCode), 0o644)

			data, changed, err := tt.generator.IncrementalUpdate(tt.filepath, []byte(tt.generatedCode))
			require.NoError(t, err)
			assert.Equal(t, tt.changed, changed)
			assert.Equal(t, tt.wantedCode, string(data))
		})
	}
}
