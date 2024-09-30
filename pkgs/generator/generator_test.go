package generator_test

import (
	"errors"
	"testing"

	"github.com/mhersson/vectorsigma/pkgs/generator"
	"github.com/mhersson/vectorsigma/pkgs/shell/mock_shell"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/spf13/afero"
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
			FS:    afero.NewMemMapFs(),
		}
	})

	ginkgo.AfterEach(func() {
		ctrl.Finish()
	})

	ginkgo.Describe("ExecuteTemplate", func() {
		ginkgo.It("should return an error if the template file does not exist", func() {
			_, err := gen.ExecuteTemplate("nonexistent.tmpl")
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("WriteFile", func() {
		ginkgo.It("should write a file successfully", func() {
			gomega.Expect(gen.WriteFile("somefile", []byte("1234567"))).To(gomega.Succeed())
			content, err := afero.ReadFile(gen.FS, "somefile")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(content).To(gomega.Equal([]byte("1234567")))
		})
	})

	ginkgo.Describe("InitializeModule", func() {
		ginkgo.It("should initialize a new go module successfully", func() {
			mockShell.EXPECT().NewCommand("go", "mod", "init", "test-module").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(nil)

			gen.Module = "test-module"
			gomega.Expect(gen.InitializeModule()).To(gomega.Succeed())
		})

		ginkgo.It("should return an error if initializing a new go module fails", func() {
			mockShell.EXPECT().NewCommand("go", "mod", "init", "").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(errors.New("oops"))

			gomega.Expect(gen.InitializeModule()).ToNot(gomega.Succeed())
		})
	})

	ginkgo.Describe("FormatCode", func() {
		ginkgo.It("should format the generated code successfully", func() {
			mockShell.EXPECT().NewCommand("go", "fmt", "path").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(nil)

			gomega.Expect(gen.FormatCode("path")).To(gomega.Succeed())
		})

		ginkgo.It("should return an error if formatting the generteated code fails", func() {
			mockShell.EXPECT().NewCommand("go", "fmt", "path").Return(mockCmd)
			mockCmd.EXPECT().Run().Return(errors.New("oops"))

			gomega.Expect(gen.FormatCode("path")).ToNot(gomega.Succeed())
		})
	})
})
