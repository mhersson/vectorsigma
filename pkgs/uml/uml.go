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
package uml

import (
	"regexp"
	"slices"
	"strings"
)

const (
	InitialState = "InitialState"
	FinalState   = "FinalState"
	// [*] --> InitialState. Before we replace the [*] with InitialState.
	firstInitialStatePattern = `^\s*\[\*\].*$`
	titlePattern             = `^title\s(.*)$`
	// InitialState --> StartingConversation.
	initialStatePattern = `^\s*` + InitialState + `\s*-->\s*(\w+)$`
	// StartingConversation: do / StartConversation(param).
	actionPattern = `^\s*(\w+):\s*(do\s*\/\s*)?(\w+)(\((.*)\))?$`
	// StartingConversation --> FinalState: [ isError ].
	guardedTransitionPattern = `^\s*(\w+)\s*-->\s*(\w+):\s*\[?\s*(\w+)\s*\]?\s*?(::\s*(\w+)(\((.*)\))?)?$`
	// StartingConversation --> FinalState.
	defaultTransitionPattern = `^\s*(\w+)\s*-->\s*(\w+)$`
	// CompositeState: state compositestate {.
	compositeStateStartPattern = `^\s*state\s*(\w+)\s*{$`

	compositeStateEndPattern = `^\s*}$`
)

type State struct {
	Name        string
	Actions     []Action
	Transitions []Transition
	Composite   Composite
}

type Composite struct {
	InitialState string
	States       map[string]*State
}

type Transition struct {
	Target string
	Guard  string
	Action *Action
}

type Action struct {
	Name   string
	Params string
}

type FSM struct {
	States       map[string]*State
	Title        string
	InitialState string
	ActionNames  []string
	GuardNames   []string
	AllStates    []string
}

func (f *FSM) Action(action string) {
	if !slices.Contains(f.ActionNames, action) {
		f.ActionNames = append(f.ActionNames, action)
	}
}

func (f *FSM) Guard(guard string) {
	if !slices.Contains(f.GuardNames, guard) {
		f.GuardNames = append(f.GuardNames, guard)
	}
}

func (f *FSM) IsTitle(line string) bool {
	re := regexp.MustCompile(titlePattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		f.Title = strings.ReplaceAll(m[1], " ", "")

		return true
	}

	return false
}

func (f *FSM) IsInitialState(line string) bool {
	re := regexp.MustCompile(initialStatePattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		if _, ok := f.States[InitialState]; !ok {
			f.States[InitialState] = &State{
				Name: InitialState,
				Transitions: []Transition{
					{
						Target: m[1],
					},
				},
			}
		} else {
			newTransition := Transition{Target: m[1]}
			f.States[InitialState].Transitions = append(f.States[InitialState].Transitions, newTransition)
		}

		return true
	}

	return false
}

func (f *FSM) IsAction(line string) bool {
	re := regexp.MustCompile(actionPattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		state := m[1]
		action := Action{}
		action.Name = m[3]

		if len(m) == 6 && m[5] != "" {
			params := strings.ReplaceAll(m[5], `"`, ``)
			action.Params = `"` + strings.ReplaceAll(params, ",", `","`) + `"`
		}

		f.Action(action.Name)

		if _, ok := f.States[state]; !ok {
			f.States[state] = &State{
				Name:    state,
				Actions: []Action{action},
			}
		} else {
			f.States[state].Actions = append(f.States[state].Actions, action)
		}

		return true
	}

	return false
}

func (f *FSM) IsGuardedTransition(line string) bool {
	re := regexp.MustCompile(guardedTransitionPattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		state := m[1]
		transition := m[2]
		guard := m[3]
		f.Guard(guard)

		var action *Action
		// Check if there is an action behind the guard.
		if len(m) >= 6 && m[5] != "" {
			action = &Action{}
			action.Name = m[5]
			if len(m) == 8 && m[7] != "" {
				action.Params = `"` + strings.ReplaceAll(m[7], ",", `","`) + `"`
				f.Action(action.Name)
			}
		}

		if _, ok := f.States[state]; !ok {
			f.States[state] = &State{
				Name: state,
				Transitions: []Transition{
					{
						Target: transition,
						Guard:  guard,
						Action: action,
					},
				},
			}
		} else {
			newTransition := Transition{Target: transition, Guard: guard, Action: action}
			f.States[state].Transitions = append(f.States[state].Transitions, newTransition)
		}

		// Make sure the target state exists.
		if _, ok := f.States[transition]; !ok {
			f.States[transition] = &State{
				Name: transition,
			}
		}

		return true
	}

	return false
}

func (f *FSM) IsDefaultTransition(line string) bool {
	re := regexp.MustCompile(defaultTransitionPattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		state := m[1]
		transition := m[2]

		if _, ok := f.States[state]; !ok {
			f.States[state] = &State{
				Name: state,
				Transitions: []Transition{
					{
						Target: transition,
					},
				},
			}
		} else {
			newTransition := Transition{Target: transition}
			f.States[state].Transitions = append(f.States[state].Transitions, newTransition)
		}

		// Make sure the target state exists.
		if _, ok := f.States[transition]; !ok {
			f.States[transition] = &State{
				Name: transition,
			}
		}

		return true
	}

	return false
}

func (f *FSM) IsCompositeStateStart(ind int, line, data string) (int, bool) {
	re := regexp.MustCompile(compositeStateStartPattern)

	m := re.FindStringSubmatch(line)
	if m != nil {
		state := m[1]
		lines := strings.Split(data, "\n")
		start := ind + 1
		end := 0
		for i := ind + 1; i < len(lines); i++ {
			if f.IsCompositeStateEnd(lines[i]) {
				end = i

				break
			}
		}

		if end == 0 {
			return end, false
		}

		compState := Parse(strings.Join(lines[start:end], "\n"))

		f.States[state] = &State{
			Name: state,
			Composite: Composite{
				InitialState: compState.InitialState,
				States:       compState.States,
			},
		}

		for k := range compState.States {
			if !slices.Contains(f.AllStates, k) {
				f.AllStates = append(f.AllStates, k)
			}
		}

		slices.Sort(f.AllStates)

		for _, v := range compState.ActionNames {
			if !slices.Contains(f.ActionNames, v) {
				f.ActionNames = append(f.ActionNames, v)
			}
		}

		for _, v := range compState.GuardNames {
			if !slices.Contains(f.GuardNames, v) {
				f.GuardNames = append(f.GuardNames, v)
			}
		}

		return end - start, true
	}

	return 0, false
}

func (f *FSM) IsCompositeStateEnd(line string) bool {
	re := regexp.MustCompile(compositeStateEndPattern)

	m := re.FindStringSubmatch(line)

	return m != nil
}

func Parse(data string) *FSM {
	fsm := new(FSM)
	fsm.States = make(map[string]*State)
	fsm.InitialState = InitialState
	fsm.AllStates = []string{}

	lines := strings.Split(data, "\n")

	for ind := 0; ind < len(lines); ind++ {
		lines[ind] = strings.ReplaceAll(lines[ind], "[dotted]", "")
		lines[ind] = strings.ReplaceAll(lines[ind], "[bold]", "")

		if strings.Contains(lines[ind], "[*]") {
			re := regexp.MustCompile(firstInitialStatePattern)
			m := re.FindStringSubmatch(lines[ind])
			if m != nil {
				lines[ind] = strings.ReplaceAll(lines[ind], "[*]", InitialState)
			} else {
				lines[ind] = strings.ReplaceAll(lines[ind], "[*]", FinalState)
			}
		}

		if fsm.IsTitle(lines[ind]) {
			continue
		}

		if fsm.IsInitialState(lines[ind]) {
			continue
		}

		if fsm.IsAction(lines[ind]) {
			continue
		}

		if fsm.IsGuardedTransition(lines[ind]) {
			continue
		}

		if fsm.IsDefaultTransition(lines[ind]) {
			continue
		}

		if i, ok := fsm.IsCompositeStateStart(ind, lines[ind], data); ok {
			ind += i

			continue
		}
	}

	for k := range fsm.States {
		if !slices.Contains(fsm.AllStates, k) {
			fsm.AllStates = append(fsm.AllStates, k)
		}
	}

	slices.Sort(fsm.AllStates)
	slices.Sort(fsm.ActionNames)
	slices.Sort(fsm.GuardNames)

	return fsm
}
