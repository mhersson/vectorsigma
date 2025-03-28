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
			name: "Ok",
			args: args{line: "InitialState --> CurrentState"}, expect: "CurrentState", want: true,
		},
		{
			name: "Ok no spaces",
			args: args{line: "InitialState-->CurrentState"}, expect: "CurrentState", want: true,
		},
		{
			name: "Not Ok",
			args: args{line: "State --> State2"}, expect: "", want: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{States: make(map[string]*uml.State)}
			if got := f.IsInitialState(tt.args.line); got != tt.want {
				t.Errorf("FSM.IsInitialState() = %v, want %v", got, tt.want)
			} else if tt.want {
				assert.Equal(t, f.States[uml.InitialState].Transitions[0].Target, tt.expect)
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
		expect uml.Action
	}{
		{
			name: "Ok", args: args{line: "State: do / action"}, expect: uml.Action{Name: "action", Params: ""},
			want: true,
		},
		{
			name: "Ok no spaces", args: args{line: "State:do/action"}, expect: uml.Action{Name: "action", Params: ""},
			want: true,
		},
		{
			name: "Not OK", args: args{line: "State --> State2: guard"}, expect: uml.Action{},
			want: false,
		},
		{
			name: "Ok params", args: args{line: "State: do / action(param1,param2)"}, expect: uml.Action{Name: "action", Params: `"param1","param2"`},
			want: true,
		},
		{
			name: "Ok params with spaces", args: args{line: "State: do / action(param1,this is a message)"}, expect: uml.Action{Name: "action", Params: `"param1","this is a message"`},
			want: true,
		},
		{
			name: "Ok params with spaces and leading space", args: args{line: "State: do / action(param1, this is a message)"}, expect: uml.Action{Name: "action", Params: `"param1","this is a message"`},
			want: true,
		},
		{
			name: "Ok params with leading spaces", args: args{line: "State: do / action(param1,     param2,param3, param4)"}, expect: uml.Action{Name: "action", Params: `"param1","param2","param3","param4"`},
			want: true,
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
				assert.Contains(t, f.ActionNames, tt.expect.Name)
				if len(tt.expect.Params) > 0 {
					assert.Equal(t, tt.expect.Params, f.States["State"].Actions[0].Params)
				}
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
		{
			name: "Ok guarded action", args: args{line: "state --> state2: [ guard ] :: myaction"}, expect: "guard",
			want: true,
		},
		{
			name: "Ok guarded action with params", args: args{line: "state --> state2: [ guard ] :: myaction(param,param1)"}, expect: "guard",
			want: true,
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

func TestFSM_IsCompositeStateStart(t *testing.T) {
	type fields struct {
		States       map[string]*uml.State
		Title        string
		InitialState string
		ActionNames  []string
		GuardNames   []string
	}
	type args struct {
		ind  int
		line string
		data string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		expect *uml.FSM
	}{
		{
			name: "Composite State Start",
			args: args{
				ind:  2,
				line: "state CompositeState {",
				data: `

[*] --> CompositeState

state CompositeState {
	[*] --> SubState1
	SubState1: do / action1
	SubState1 --> SubState2

	SubState2 --> [*]
}
CompositeState --> [*]
`,
			},
			want: true,
			expect: &uml.FSM{
				AllStates:   []string{"FinalState", "InitialState", "SubState1", "SubState2"},
				ActionNames: []string{"action1"},
				States: map[string]*uml.State{
					"CompositeState": {
						Name: "CompositeState",
						Composite: uml.Composite{
							InitialState: uml.InitialState,
							States: map[string]*uml.State{
								"InitialState": {
									Name: "InitialState",
									Transitions: []uml.Transition{
										{Target: "SubState1"},
									},
								},
								"SubState1": {
									Name: "SubState1",
									Actions: []uml.Action{
										{Name: "action1", Params: ""},
									},
									Transitions: []uml.Transition{
										{Target: "SubState2"},
									},
								},
								"SubState2": {
									Name: "SubState2",
									Transitions: []uml.Transition{
										{Target: "FinalState"},
									},
								},
								uml.FinalState: {
									Name: uml.FinalState,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Not Composite State Start",
			args: args{
				ind:  0,
				line: "state NotCompositeState",
				data: `
state NotCompositeState
[*] --> SubState1
SubState1 --> SubState2
`,
			},
			want:   false,
			expect: &uml.FSM{States: map[string]*uml.State{}},
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &uml.FSM{States: make(map[string]*uml.State)}
			if _, got := f.IsCompositeStateStart(tt.args.ind, tt.args.line, tt.args.data); got != tt.want {
				t.Errorf("FSM.IsCompositeStateStart() = %v, want %v", got, tt.want)
			} else {
				assert.Equal(t, tt.expect, f)
			}
		})
	}
}

func TestParseSimple(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *uml.FSM
	}{
		{
			name: "Simple",
			want: &uml.FSM{
				InitialState: uml.InitialState,
				States: map[string]*uml.State{
					uml.InitialState: {
						Name: uml.InitialState,
						Transitions: []uml.Transition{
							{Target: "Red", Guard: ""},
						},
					},
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
							{Target: "Yellow", Guard: ""},
						},
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
					uml.FinalState: {
						Name: uml.FinalState,
					},
				},
				Title:       "TrafficLight",
				ActionNames: []string{"SwitchIn"},
				GuardNames:  []string{"IsError", "NotGonnaHappen"},
				AllStates: []string{
					"FinalState", "FlashingYellow", "Green", "InitialState", "Red", "Yellow",
				},
			},
			args: args{
				data: `
@startuml

title Traffic Light
[*] -down-> Red
Red: do / SwitchIn(5)
Red -[dotted]-> [*]: [ IsError ]
Red --> Green: [ NotGonnaHappen ]
Red -[bold]-> Yellow

Yellow: do / SwitchIn(1)
Yellow -[dotted]left-> [*]: [ IsError]
Yellow -right-> Green

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

func TestParseComposite(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *uml.FSM
	}{
		{
			name: "Composite State",
			want: &uml.FSM{
				InitialState: uml.InitialState,
				States: map[string]*uml.State{
					uml.InitialState: {
						Name: uml.InitialState,
						Transitions: []uml.Transition{
							{Target: "Red", Guard: ""},
						},
					},
					"Red": {
						Name: "Red",
						Transitions: []uml.Transition{
							{
								Target: "FinalState",
								Guard:  "IsError",
							},
							{Target: "Yellow", Guard: ""},
						},
						Composite: uml.Composite{
							InitialState: uml.InitialState,
							States: map[string]*uml.State{
								uml.InitialState: {
									Name: uml.InitialState,
									Transitions: []uml.Transition{
										{Target: "InRed", Guard: ""},
									},
								},
								"InRed": {
									Name: "InRed",
									Actions: []uml.Action{{
										Name:   "SwitchIn",
										Params: `"5"`,
									}},
									Transitions: []uml.Transition{
										{
											Target: "FinalState",
											Guard:  "IsError",
										},
										{Target: uml.FinalState, Guard: ""},
									},
								},
								uml.FinalState: {
									Name: uml.FinalState,
								},
							},
						},
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
					uml.FinalState: {
						Name: uml.FinalState,
					},
				},
				Title:       "TrafficLight",
				ActionNames: []string{"SwitchIn"},
				GuardNames:  []string{"IsError"},
				AllStates: []string{
					"FinalState", "FlashingYellow", "Green", "InRed", "InitialState", "Red", "Yellow",
				},
			},
			args: args{
				data: `
@startuml

title Traffic Light
[*] --> Red

state Red {
	[*] --> InRed
	InRed: do / SwitchIn(5)
	InRed --> [*]: [ IsError ]
	InRed --> [*]
}
Red -[dotted]-> [*]: [ IsError ]
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

func TestParseGuardedAction(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want *uml.FSM
	}{
		{
			name: "Simple",
			want: &uml.FSM{
				InitialState: uml.InitialState,
				States: map[string]*uml.State{
					uml.InitialState: {
						Name: uml.InitialState,
						Transitions: []uml.Transition{
							{Target: "Red", Guard: ""},
						},
					},
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
								Action: &uml.Action{
									Name:   "AlwaysGreen",
									Params: `"force1","force2"`,
								},
							},
							{Target: "Yellow", Guard: ""},
						},
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
						}, {Target: "FinalState", Guard: ""}},
					},
					"Green": {
						Name: "Green",
						Actions: []uml.Action{{
							Name:   "SwitchIn",
							Params: `"5"`,
						}},
						Transitions: []uml.Transition{
							{Target: "FinalState", Guard: ""},
						},
					},
					uml.FinalState: {
						Name: uml.FinalState,
					},
				},
				Title:       "TrafficLight",
				ActionNames: []string{"AlwaysGreen", "SwitchIn"},
				GuardNames:  []string{"IsError", "NotGonnaHappen"},
				AllStates: []string{
					"FinalState", "Green", "InitialState", "Red", "Yellow",
				},
			},
			args: args{
				data: `
@startuml

title Traffic Light
[*] --> Red
Red: do / SwitchIn(5)
Red -[dotted]-> [*]: [ IsError ]
Red --> Green: NotGonnaHappen::AlwaysGreen(force1,force2)
Red --> Yellow

Yellow: do / SwitchIn(1)
Yellow -[dotted]-> [*]: [ IsError]
Yellow --> [*]

Green: do / SwitchIn(5)
Green --> [*]

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
