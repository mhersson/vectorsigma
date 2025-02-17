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
	"github.com/mhersson/vectorsigma/internal/statemachine"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	update      = flag.Bool("update", false, "update golden files")
	vectorsigma *cobra.Command
)

// Just use the standard go test.
// nolint: paralleltest
func Test_IntegrationTest(t *testing.T) {
	tests := []struct {
		name           string
		testdatafolder string
		init           bool
		input          string
		output         string
		pkg            string
		// arguments      []string
	}{
		{
			name:           "Initialize module",
			testdatafolder: "new_module",
			output:         "output",
			init:           true,
			input:          "../uml/traffic-lights.plantuml",
			pkg:            "fsm",
		},
		{
			name:           "Generate package",
			testdatafolder: "package",
			output:         "output",
			init:           false,
			input:          "../uml/traffic-lights.plantuml",
			pkg:            "fsm",
		},
	}

	rootDir, _ := os.Getwd()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := filepath.Join("testdata", tt.testdatafolder)
			outputPath := filepath.Join(testPath, tt.output)
			goldenPath := filepath.Join(testPath, "golden")

			err := cleanup(outputPath)
			require.NoError(t, err)

			err = os.Chdir(testPath)
			require.NoError(t, err)

			vectorsigma = cmd.RootCmd
			vectorsigma.SetArgs([]string{"generate", "--input", "weoverridethis"})

			cmd.SM = statemachine.New()
			cmd.SM.ExtendedState.Init = tt.init
			cmd.SM.ExtendedState.Input = tt.input
			cmd.SM.ExtendedState.Package = tt.pkg
			cmd.SM.ExtendedState.Output = tt.output

			err = vectorsigma.Execute()
			require.NoError(t, err)

			err = os.Chdir(rootDir)
			require.NoError(t, err)

			if *update {
				err = os.RemoveAll(goldenPath)
				require.NoError(t, err)

				err = copyFiles(outputPath, goldenPath)
				require.NoError(t, err)
			}

			// Check the output
			err = checkOutput(t, goldenPath, outputPath)
			require.NoError(t, err)
		})
	}
}

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
func checkOutput(t *testing.T, goldenPath, outputPath string) error {
	return filepath.WalkDir(goldenPath, func(path string, d fs.DirEntry, err error) error {
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
		require.NoError(t, err)

		assert.Equal(t, genstr, string(golden), "%s: %s", outputPath, path)

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
