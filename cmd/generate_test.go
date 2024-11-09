package cmd_test

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mhersson/vectorsigma/cmd"
	"github.com/spf13/cobra"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	t.Parallel()
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Integration tests")
}

var (
	update      = flag.Bool("update", false, "update golden files")
	initCmd     = []string{"generate", "-i", "../uml/traffic-lights.plantuml", "--init"}
	genCmd      = []string{"generate", "-i", "../uml/traffic-lights.plantuml"}
	vectorsigma *cobra.Command
)

var _ = ginkgo.DescribeTable("Generate Integration tests", ginkgo.Label("integration"),
	func(directory string, arguments []string) {
		cmd.InputParams.Init = false
		vectorsigma = cmd.RootCmd

		testPath := filepath.Join("testdata", directory)
		outputPath := filepath.Join(testPath, "output")
		goldenPath := filepath.Join(testPath, "golden")
		rootDir, _ := os.Getwd()

		err := cleanup(outputPath)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		err = os.Chdir(testPath)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		// Add relative output path to arguments
		arguments = append(arguments, "--output", "output")

		vectorsigma.SetArgs(arguments)

		err = vectorsigma.Execute()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		err = os.Chdir(rootDir)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		// If update is set replace the all expected files with the generated ones
		if *update {
			err = os.RemoveAll(goldenPath)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			err = copyFiles(outputPath, goldenPath)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		}

		// Check the output
		err = checkOutput(goldenPath, outputPath)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	},

	ginkgo.Entry("Initialize new module", "initializeNewModule", initCmd),
	ginkgo.Entry("Generate fsm package", "generateBasicFSM", genCmd),
)

// checkOutput compares the generated output files with the expected golden files.
//
// It walks through the expectedPath directory and for each file, it reads the corresponding
// file from the outputPath directory. It then compares the content of the generated file
// with the expected content. If the update flag is set, it updates the golden files with
// the generated content.
//
// Parameters:
// - goldenPath: The directory containing the expected golden files.
// - outputPath: The directory containing the generated output files.
//
// Returns:
// - An error if any file comparison fails or if there are issues reading the files.
//
//nolint:wrapcheck
func checkOutput(goldenPath, outputPath string) error {
	return filepath.WalkDir(goldenPath, func(path string, d fs.DirEntry, err error) error {
		ginkgo.By("Checking " + path)

		if err != nil {
			return err
		}

		if d.Name() == "go.mod" {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		filepath := strings.Replace(path, goldenPath, outputPath, 1)

		var genBytes []byte

		genBytes, err = os.ReadFile(filepath)
		if err != nil {
			return err //nolint:wrapcheck
		}

		genstr := string(genBytes)
		golden, err := os.ReadFile(path)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		gomega.Expect(genstr).To(gomega.Equal(string(golden)), fmt.Sprintf("%s: %s", outputPath, path))

		return nil
	})
}

// cleanup removes any existing files at the outputPath and creates the
// necessary directories. Returns an error if any operation fails.
//
//nolint:wrapcheck
func cleanup(outputPath string) error {
	// Remove potential leftovers
	err := os.RemoveAll(outputPath)
	if err != nil {
		return err
	}

	// Copy wanted existing file to the destination folder
	err = os.MkdirAll(outputPath, 0o755)
	if err != nil {
		return err
	}

	return nil
}

func copyFiles(inputPath, outputPath string) error {
	files, err := os.ReadDir(inputPath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			if err = copyFiles(filepath.Join(inputPath, f.Name()), filepath.Join(outputPath, f.Name())); err != nil {
				return err
			}
		} else {
			if err = copyFile(filepath.Join(inputPath, f.Name()), filepath.Join(outputPath, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	if err := createDirIfNotExists(dst); err != nil {
		return err
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", src, err)
	}

	if err = os.WriteFile(dst, data, 0o644); err != nil { //nolint:gosec
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}

	return nil
}

func createDirIfNotExists(directory string) error {
	dir := filepath.Dir(directory)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return nil
}
