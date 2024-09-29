package generator_test

import (
	"testing"

	"github.com/mhersson/vectorsigma/pkgs/generator"
	"github.com/mhersson/vectorsigma/pkgs/shell/mock_shell"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestGenerator(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Generator")
}

var _ = ginkgo.Describe("Generator", func() {
	var (
		ctrl      *gomock.Controller
		gen       *generator.Generator
		mockShell *mock_shell.MockInterface
		mockCmd   *mock_shell.MockCmdRunner
	)

	ginkgo.BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		mockShell = mock_shell.NewMockInterface(ctrl)
		mockCmd = mock_shell.NewMockCmdRunner(ctrl)
		gen = &generator.Generator{
			Shell: mockShell,
		}
	})

	ginkgo.AfterEach(func() {
		ctrl.Finish()
	})

	ginkgo.Describe("InitializeModule", func() {
		ginkgo.It("should initialize a new go module successfully", func() {
			mockShell.EXPECT().NewCommand("go", "mod", "init", "test-module").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(nil)

			gen.Module = "test-module"
			gomega.Expect(gen.InitializeModule()).To(gomega.Succeed())
		})
	})
})
