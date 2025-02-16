package statemachine_test

import (
	"path/filepath"
	"testing"

	"github.com/mhersson/vectorsigma/internal/statemachine"
	"github.com/mhersson/vectorsigma/pkgs/generator"
	"github.com/mhersson/vectorsigma/pkgs/shell/mock_shell"
	"github.com/mhersson/vectorsigma/pkgs/uml"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestVectorSigma_InitializeAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK", fields: fields{
			context:       &statemachine.Context{},
			ExtendedState: &statemachine.ExtendedState{}},
			wantErr: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.InitializeAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.InitializeAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				assert.NotEqual(t, "", fsm.ExtendedState.Output)
				assert.NotEqual(t, "", fsm.ExtendedState.Module)
				assert.Equal(t, fsm.Context.Generator.Module, fsm.ExtendedState.Module)
				assert.Equal(t, fsm.Context.Generator.Package, fsm.ExtendedState.Package)
			}
		})
	}
}

func TestVectorSigma_LoadInputAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	var fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "input.md", []byte("# Markdown"), 0664)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context:       &statemachine.Context{Generator: &generator.Generator{FS: fs}},
				ExtendedState: &statemachine.ExtendedState{Input: "input.md"}},
			wantErr: false,
		},
		{name: "NOT OK",
			fields: fields{
				context:       &statemachine.Context{Generator: &generator.Generator{FS: fs}},
				ExtendedState: &statemachine.ExtendedState{Input: "invalid.md"}},
			wantErr: true},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.LoadInputAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.LoadInputAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVectorSigma_ExtractUMLAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				ExtendedState: &statemachine.ExtendedState{
					InputData: "# Markdown ```plantuml\n@startuml```",
				}},
			wantErr: false},
		{name: "Not OK - typo in markdown",
			fields: fields{
				ExtendedState: &statemachine.ExtendedState{
					InputData: "# Markdown ``plantuml\n@startuml```",
				}},
			wantErr: true},
		{name: "Not OK - missisng end of block",
			fields: fields{
				ExtendedState: &statemachine.ExtendedState{
					InputData: "# Markdown ``plantuml\n@startuml",
				}},
			wantErr: true},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.ExtractUMLAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.ExtractUMLAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				assert.Equal(t, "\n@startuml", fsm.ExtendedState.InputData)
			}
		})
	}
}

func TestVectorSigma_ParseUMLAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{}},
				ExtendedState: &statemachine.ExtendedState{
					InputData: "# Markdown ```plantuml\n@startuml\ntitle test title\n\nskin rose```",
				}},
			wantErr: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.ParseUMLAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.ParseUMLAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				assert.Equal(t, "testtitle", fsm.Context.Generator.FSM.Title)
			}
		})
	}
}

func TestVectorSigma_GenerateStateMachineAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{FSM: &uml.FSM{}, Package: "unittest"}},
				ExtendedState: &statemachine.ExtendedState{
					Package:       "unittest",
					GeneratedData: make(map[string][]byte),
				}},
			wantErr: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.GenerateStateMachineAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.GenerateStateMachineAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				for k, v := range fsm.ExtendedState.GeneratedData {
					assert.Contains(t, string(v), "package unittest", k)
				}
			}
		})
	}
}

func TestVectorSigma_CreateOutputFolderAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	var fs = afero.NewMemMapFs()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{FS: fs}},
				ExtendedState: &statemachine.ExtendedState{
					Output:  "outputfolder",
					Package: "statemachine"}},
			wantErr: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.CreateOutputFolderAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.CreateOutputFolderAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				exists, _ := afero.Exists(fs, filepath.Join(tt.fields.ExtendedState.Output, tt.fields.ExtendedState.Package))
				assert.True(t, exists)
			}
		})
	}
}

func TestVectorSigma_WriteGeneratedFilesAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	var fs = afero.NewMemMapFs()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{FS: fs}},
				ExtendedState: &statemachine.ExtendedState{
					Output:        "outputfolder/statemachine",
					Package:       "statemachine",
					GeneratedData: map[string][]byte{"actions.go": []byte("actions"), "guards.go": []byte("guards")},
				}},

			wantErr: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.WriteGeneratedFilesAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.WriteGeneratedFilesAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				assert.NoError(t, fs.MkdirAll(tt.fields.ExtendedState.Output, 0o755))

				for f, c := range tt.fields.ExtendedState.GeneratedData {
					filename := filepath.Join(tt.fields.ExtendedState.Output, f)

					exists, err := afero.Exists(fs, filename)
					assert.NoError(t, err)
					assert.True(t, exists)

					d, err := afero.ReadFile(fs, filename)
					assert.NoError(t, err)
					assert.Equal(t, c, d)
				}
			}
		})
	}
}

func TestVectorSigma_InitializeGoModuleAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	mockShell := mock_shell.NewMockInterface(t)
	mockCmd := mock_shell.NewMockCmdRunner(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{
					Shell:  mockShell,
					Module: "test-module"},
				},
			},
			wantErr: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			mockShell.EXPECT().NewCommand("go", "mod", "init", "test-module").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(nil)

			if err := fsm.InitializeGoModuleAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.InitializeGoModuleAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVectorSigma_GenerateMainFileAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{FSM: &uml.FSM{}, Package: "unittest"}},
				ExtendedState: &statemachine.ExtendedState{
					Package:       "unittest",
					GeneratedData: make(map[string][]byte),
				},
			},
			wantErr: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if err := fsm.GenerateMainFileAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.GenerateMainFileAction() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				for k, v := range fsm.ExtendedState.GeneratedData {
					assert.Equal(t, "main.go", k)
					assert.Contains(t, string(v), "package main", k)
				}
			}
		})
	}
}

func TestVectorSigma_FormatCodeAction(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	type args struct {
		params []string
	}

	mockShell := mock_shell.NewMockInterface(t)
	mockCmd := mock_shell.NewMockCmdRunner(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "OK",
			fields: fields{
				context: &statemachine.Context{Generator: &generator.Generator{Shell: mockShell}},
				ExtendedState: &statemachine.ExtendedState{
					Output: "out",
				},
			},
			wantErr: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}

			mockShell.EXPECT().NewCommand("go", "fmt", "out/...").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(nil)

			if err := fsm.FormatCodeAction(tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("VectorSigma.FormatCodeAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
