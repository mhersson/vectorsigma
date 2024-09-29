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
	"testing"

	"github.com/mhersson/vectorsigma/pkgs/uml"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestFSM(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "UML parser")
}

var _ = ginkgo.Describe("FSM", func() {
	var fsm *uml.FSM

	ginkgo.BeforeEach(func() {
		fsm = &uml.FSM{
			States: make(map[string]*uml.State),
		}
	})

	ginkgo.Describe("IsTitle", func() {
		ginkgo.It("should parse valid titles", func() {
			gomega.Expect(fsm.IsTitle("title My Title2")).To(gomega.BeTrue())
			gomega.Expect(fsm.Title).To(gomega.Equal("MyTitle2"))

			gomega.Expect(fsm.IsTitle("title Super Valid but a bit long title")).To(gomega.BeTrue())
			gomega.Expect(fsm.Title).To(gomega.Equal("SuperValidbutabitlongtitle"))
		})

		ginkgo.It("should not parse invalid titles", func() {
			gomega.Expect(fsm.IsTitle("wrong title")).To(gomega.BeFalse())
			gomega.Expect(fsm.Title).To(gomega.Equal(""))

			gomega.Expect(fsm.IsTitle("title: invalid")).To(gomega.BeFalse())
			gomega.Expect(fsm.Title).To(gomega.Equal(""))
		})
	})

	ginkgo.Describe("IsInitialState", func() {
		ginkgo.It("should parse valid initial states", func() {
			gomega.Expect(fsm.IsInitialState("[*] --> StartingConversation")).To(gomega.BeTrue())
			gomega.Expect(fsm.InitialState).To(gomega.Equal("StartingConversation"))

			gomega.Expect(fsm.IsInitialState("[*]--> StartingConversation2")).To(gomega.BeTrue())
			gomega.Expect(fsm.InitialState).To(gomega.Equal("StartingConversation2"))
		})

		ginkgo.It("should not parse invalid initial states", func() {
			gomega.Expect(fsm.IsInitialState("[* ] --> StartingConversation")).To(gomega.BeFalse())
			gomega.Expect(fsm.InitialState).To(gomega.Equal(""))

			gomega.Expect(fsm.IsInitialState(" [*]--> StartingConversation2")).To(gomega.BeFalse())
			gomega.Expect(fsm.InitialState).To(gomega.Equal(""))

			gomega.Expect(fsm.IsInitialState("[*] -> StartingConversation")).To(gomega.BeFalse())
			gomega.Expect(fsm.InitialState).To(gomega.Equal(""))
		})
	})

	ginkgo.Describe("IsAction", func() {
		ginkgo.It("should parse valid actions", func() {
			gomega.Expect(fsm.IsAction("StartingConversation: do / StartConversation")).To(gomega.BeTrue())
			gomega.Expect(fsm.ActionNames).To(gomega.ContainElement("StartConversation"))
			gomega.Expect(fsm.States["StartingConversation"].Actions).To(gomega.ContainElement(uml.Action{Name: "StartConversation"}))

			gomega.Expect(fsm.IsAction("StartingConversation:do/StartConversation")).To(gomega.BeTrue())
			gomega.Expect(fsm.ActionNames).To(gomega.ContainElement("StartConversation"))
			gomega.Expect(fsm.States["StartingConversation"].Actions).To(gomega.ContainElement(uml.Action{Name: "StartConversation"}))
		})

		ginkgo.It("should parse actions with parameters", func() {
			gomega.Expect(fsm.IsAction("StartingConversation: do / StartConversation(5)")).To(gomega.BeTrue())
			gomega.Expect(fsm.ActionNames).To(gomega.ContainElement("StartConversation"))
			gomega.Expect(fsm.States["StartingConversation"].Actions).To(gomega.ContainElement(uml.Action{Name: "StartConversation", Params: `"5"`}))

			gomega.Expect(fsm.IsAction("StartingConversation:do/StartConversation(one,two)")).To(gomega.BeTrue())
			gomega.Expect(fsm.ActionNames).To(gomega.ContainElement("StartConversation"))
			gomega.Expect(fsm.States["StartingConversation"].Actions).To(gomega.ContainElement(uml.Action{Name: "StartConversation", Params: `"one","two"`}))
		})

		ginkgo.It("should not parse invalid actions", func() {
			gomega.Expect(fsm.IsAction("StartingConversation do / StartConversation")).To(gomega.BeFalse())
			gomega.Expect(fsm.ActionNames).To(gomega.BeEmpty())
			gomega.Expect(fsm.States).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("IsGuardedTransition", func() {
		ginkgo.It("should parse valid guarded transitions", func() {
			gomega.Expect(fsm.IsGuardedTransition("StartingConversation --> FinalState: [ IsError ]")).To(gomega.BeTrue())
			gomega.Expect(fsm.GuardNames).To(gomega.ContainElement("IsError"))
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.ContainElement(uml.Transition{Target: "FinalState", Guard: "IsError"}))
			gomega.Expect(fsm.States["FinalState"]).NotTo(gomega.BeNil())

			gomega.Expect(fsm.IsGuardedTransition("StartingConversation--> FinalState:[IsError ]")).To(gomega.BeTrue())
			gomega.Expect(fsm.GuardNames).To(gomega.ContainElement("IsError"))
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.ContainElement(uml.Transition{Target: "FinalState", Guard: "IsError"}))
			gomega.Expect(fsm.States["FinalState"]).NotTo(gomega.BeNil())
		})

		ginkgo.It("should not parse invalid guarded transitions", func() {
			gomega.Expect(fsm.IsGuardedTransition("StartingConversation --> FinalState: IsError")).To(gomega.BeFalse())
			gomega.Expect(fsm.GuardNames).To(gomega.BeEmpty())
			gomega.Expect(fsm.States).To(gomega.BeEmpty())
		})

		ginkgo.It("should parse valid guarded transitions to already parsed state", func() {
			fsm.States["StartingConversation"] = &uml.State{Name: "StartingConversation"}

			gomega.Expect(fsm.IsGuardedTransition("StartingConversation --> FinalState: [ IsError ]")).To(gomega.BeTrue())
			gomega.Expect(fsm.GuardNames).To(gomega.ContainElement("IsError"))
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.ContainElement(uml.Transition{Target: "FinalState", Guard: "IsError"}))
			gomega.Expect(fsm.States["FinalState"]).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("IsDefaultTransition", func() {
		ginkgo.It("should parse valid default transitions", func() {
			gomega.Expect(fsm.IsDefaultTransition("StartingConversation --> FinalState")).To(gomega.BeTrue())
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.ContainElement(uml.Transition{Target: "FinalState"}))
			gomega.Expect(fsm.States["FinalState"]).NotTo(gomega.BeNil())
		})

		ginkgo.It("should not parse invalid default transitions", func() {
			fsm.States["StartingConversation"] = &uml.State{Name: "StartingConversation"}

			gomega.Expect(fsm.IsDefaultTransition("StartingConversation-> FinalState")).To(gomega.BeFalse())
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.BeNil())
		})

		ginkgo.It("should parse valid default transitions to already found state", func() {
			fsm.States["StartingConversation"] = &uml.State{Name: "StartingConversation"}

			gomega.Expect(fsm.IsDefaultTransition("StartingConversation --> FinalState")).To(gomega.BeTrue())
			gomega.Expect(fsm.States["StartingConversation"].Transitions).To(gomega.ContainElement(uml.Transition{Target: "FinalState"}))
			gomega.Expect(fsm.States["FinalState"]).NotTo(gomega.BeNil())
		})
	})
})
