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
	titlePattern = `^title\s(.*)$`
	// FinalState --> StartingConversation.
	initialStatePattern = `^\[\*\]\s*-->\s*(\w+)$`
	// StartingConversation: do / StartConversation(param).
	actionPattern = `^(\w+):\s*do\s*\/\s*(\w+)(\((.*)\))?$`
	// StartingConversation --> FinalState: [ isError ].
	guardedTransitionPattern = `^(\w+)\s*-->\s*(\w+):\s*\[?\s*(\w+)\s*\]?$`
	// StartingConversation --> FinalState.
	defaultTransitionPattern = `^(\w+)\s*-->\s*(\w+)$`
	FinalState               = "FinalState"
)

type State struct {
	Name        string
	Actions     []Action
	Transitions []Transition
}

type Transition struct {
	Target string
	Guard  string
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
		f.InitialState = m[1]

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
		action.Name = m[2]

		if len(m) == 5 && m[4] != "" {
			params := strings.ReplaceAll(m[4], `"`, ``)
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

		if _, ok := f.States[state]; !ok {
			f.States[state] = &State{
				Name: state,
				Transitions: []Transition{
					{
						Target: transition,
						Guard:  guard,
					},
				},
			}
		} else {
			newTransition := Transition{Target: transition, Guard: guard}
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

func Parse(data string) *FSM {
	fsm := new(FSM)
	fsm.States = make(map[string]*State)

	for _, line := range strings.Split(data, "\n") {
		line = strings.ReplaceAll(line, "[dotted]", "")
		line = strings.ReplaceAll(line, "[bold]", "")

		if !strings.HasPrefix(line, "[*]") {
			line = strings.ReplaceAll(line, "[*]", FinalState)
		}

		if fsm.IsTitle(line) {
			continue
		}

		if fsm.IsInitialState(line) {
			continue
		}

		if fsm.IsAction(line) {
			continue
		}

		if fsm.IsGuardedTransition(line) {
			continue
		}

		if fsm.IsDefaultTransition(line) {
			continue
		}
	}

	slices.Sort(fsm.ActionNames)
	slices.Sort(fsm.GuardNames)

	return fsm
}
