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
package uml_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mhersson/vectorsigma/pkgs/uml"
	"github.com/stretchr/testify/assert"
)

func TestFSM_IsTitle(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect string
	}{
		{
			name: "Ok", args: args{line: "title My Title"}, expect: "MyTitle", want: true,
		},
		{

			name: "Not Ok", args: args{line: "No Title"}, expect: "", want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{}
			if got := f.IsTitle(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsTitle() = %v, want %v", got, tt.want)
			} else {
				assert.Equal(t, f.Title, tt.expect)
			}
		})
	}
}

func TestFSM_IsInitialState(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect string
	}{
		{
			name: "Ok", args: args{line: "[*] --> InitialState"}, expect: "InitialState", want: true,
		},
		{
			name: "Ok no spaces", args: args{line: "[*]-->InitialState"}, expect: "InitialState", want: true,
		},
		{

			name: "Not Ok", args: args{line: "State --> State2"}, expect: "", want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{}
			if got := f.IsInitialState(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsInitialState() = %v, want %v", got, tt.want)
			} else {
				assert.Equal(t, f.InitialState, tt.expect)
			}
		})
	}
}

func TestFSM_IsAction(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect string
	}{
		{
			name: "Ok", args: args{line: "State: do / action"}, expect: "action",
			want: true,
		},
		{
			name: "Ok no spaces", args: args{line: "State:do/action"}, expect: "action",
			want: true,
		},
		{
			name: "Not OK", args: args{line: "State --> State2: guard"}, expect: "",
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{ActionNames: []string{}, States: make(map[string]*uml.State)}
			if got := f.IsAction(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsAction() = %v, want %v", got, tt.want)
			} else if tt.want {
				assert.Contains(t, f.ActionNames, tt.expect)
			}
		})
	}
}

func TestFSM_IsGuardedTransition(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect string
	}{
		{
			name: "Ok", args: args{line: "state --> state2: [ guard ]"}, expect: "guard",
			want: true,
		},
		{
			name: "Ok no spaces", args: args{line: "state-->state2:[guard]"}, expect: "guard",
			want: true,
		},
		{
			name: "Not Ok", args: args{line: "State: do / action"}, expect: "",
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{GuardNames: []string{}, States: make(map[string]*uml.State)}

			if got := f.IsGuardedTransition(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsGuardedTransition() = %v, want %v", got, tt.want)
			} else if tt.want {
				assert.Contains(t, f.GuardNames, tt.expect)
			}
		})
	}
}

func TestFSM_IsDefaultTransition(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect string
	}{
		{
			name: "Ok", args: args{line: "state --> state2"}, expect: "state2",
			want: true,
		},
		{
			name: "Ok no spaces", args: args{line: "state-->state2"}, expect: "state2",
			want: true,
		},
		{
			name: "Ok", args: args{line: "state --> state2: guard"}, expect: "",
			want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{States: make(map[string]*uml.State)}
			if got := f.IsDefaultTransition(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsDefaultTransition() = %v, want %v", got, tt.want)
			} else if tt.want {
				assert.Contains(t, f.States["state"].Transitions, uml.Transition{Target: tt.expect})
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *uml.FSM
	}{
		{
			want: &uml.FSM{
				States: map[string]*uml.State{
					"Red": {
						Name: "Red",
						Actions: []uml.Action{{
							Name:   "SwitchIn",
							Params: `"5"`,
						}},
						Transitions: []uml.Transition{
							{
								Target: "FinalState",
								Guard:  "IsError",
							},
							{
								Target: "Green",
								Guard:  "NotGonnaHappen",
							},
							{Target: "Yellow", Guard: ""}},
					},
					"Yellow": {
						Name: "Yellow",
						Actions: []uml.Action{{
							Name:   "SwitchIn",
							Params: `"1"`,
						}},
						Transitions: []uml.Transition{{
							Target: "FinalState",
							Guard:  "IsError",
						}, {Target: "Green", Guard: ""}},
					},
					"FlashingYellow": {
						Name: "FlashingYellow",
						Actions: []uml.Action{{
							Name:   "SwitchIn",
							Params: `"3"`,
						}},
						Transitions: []uml.Transition{{
							Target: "FinalState",
							Guard:  "IsError",
						}, {Target: "Red", Guard: ""}},
					},
					"Green": {
						Name: "Green",
						Actions: []uml.Action{{
							Name:   "SwitchIn",
							Params: `"5"`,
						}},
						Transitions: []uml.Transition{{
							Target: "FinalState",
							Guard:  "IsError",
						}, {Target: "FlashingYellow", Guard: ""}},
					},
					"FinalState": {
						Name: "FinalState",
					},
				},
				Title:        "TrafficLight",
				InitialState: "Red",
				ActionNames:  []string{"SwitchIn"},
				GuardNames:   []string{"IsError", "NotGonnaHappen"},
			},
			args: args{
				data: `
@startuml

title Traffic Light
[*] --> Red
Red: do / SwitchIn(5)
Red -[dotted]-> [*]: [ IsError ]
Red --> Green: [ NotGonnaHappen ]
Red --> Yellow

Yellow: do / SwitchIn(1)
Yellow -[dotted]-> [*]: [ IsError]
Yellow --> Green

FlashingYellow: do / SwitchIn(3)
FlashingYellow -[dotted]-> [*]: [ IsError ]
FlashingYellow -[bold]-> Red

Green: do / SwitchIn(5)
Green -[dotted]-> [*]: [ IsError ]
Green --> FlashingYellow

@enduml
`,
			},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := uml.Parse(tt.args.data); !cmp.Equal(got, tt.want) {
				fmt.Println(cmp.Diff(got, tt.want))
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
