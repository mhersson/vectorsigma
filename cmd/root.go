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
package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mhersson/vectorsigma/internal/statemachine"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)

const (
	apiKindFlag    = "api-kind"
	apiVersionFlag = "api-version"
	groupFlag      = "group"
	initFlag       = "init"
	inputFlag      = "input"
	moduleFlag     = "module"
	operatorFlag   = "operator"
	outputFlag     = "output"
	packageFlag    = "package"
)

var SM *statemachine.VectorSigma

var RootCmd = &cobra.Command{
	Use:   "vectorsigma",
	Short: "VectorSigma is a Finite State Machine generator",
	Long: `VectorSigma is a Finite State Machine (FSM) generator.

VectorSigma takes PlantUML as input and generates a FSM.
`,
	Version: getVersionInfo(),
	PreRun: func(cmd *cobra.Command, _ []string) {
		if operator, _ := cmd.Flags().GetBool("operator"); operator {
			_ = cmd.MarkFlagRequired("api-kind")
			_ = cmd.MarkFlagRequired("api-version")
		}
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		SM.ExtendedState.VectorSigmaVersion = cmd.Version
		if SM.ExtendedState.Module == "" {
			SM.ExtendedState.Module = getModuleName()
		}

		return SM.Run()
	},
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new module",
	Long: `Initialize a new a new Go module with a FSM
application generated from your UML diagram.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		SM.ExtendedState.VectorSigmaVersion = getVersionInfo()
		SM.ExtendedState.Init = true

		return SM.Run()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getVersionInfo() string {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					CommitSHA = setting.Value[:8]
				}
				if setting.Key == "vcs.time" {
					BuildTime = setting.Value
				}
			}

			Version = info.Main.Version
			if Version == "(devel)" {
				return Version
			}
		}
	}

	if CommitSHA == "unknown" || BuildTime == "unknown" {
		return Version
	}

	return fmt.Sprintf("%s (commit: %s, built at: %s)", Version, CommitSHA, BuildTime)
}

func init() {
	SM = statemachine.New()

	addCommonFlags(RootCmd)
	addCommonFlags(InitCmd)

	RootCmd.AddCommand(InitCmd)
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	RootCmd.Flags().StringVarP(&SM.ExtendedState.APIKind, apiKindFlag, "k", "", "API kind (only used if generating a k8s operator)")
	RootCmd.Flags().StringVarP(&SM.ExtendedState.APIVersion, apiVersionFlag, "v", "", "API version (only used if generating a k8s operator)")
	RootCmd.Flags().StringVarP(&SM.ExtendedState.Group, groupFlag, "g", "", "Group (only used if generating a k8s operator)")
	RootCmd.Flags().BoolVarP(&SM.ExtendedState.Operator, operatorFlag, "O", false, "generate fsm for a k8s operator")
	RootCmd.Flags().StringVarP(&SM.ExtendedState.Output, outputFlag, "o", "",
		"The output path of the generated FSM (default current working directory)")
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&SM.ExtendedState.Module, moduleFlag, "m", "",
		"Name of new go module (default current directory name)")
	cmd.Flags().StringVarP(&SM.ExtendedState.Input, inputFlag, "i", "", "The UML input file")
	_ = cmd.MarkFlagRequired(inputFlag)
	cmd.Flags().StringVarP(&SM.ExtendedState.Package, packageFlag, "p", "statemachine",
		"The package name of the generated FSM")
}

func getModuleName() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		fmt.Println("Failed to read go.mod:", err)

		return ""
	}

	modFile, err := modfile.ParseLax("go.mod", content, nil)
	if err != nil {
		fmt.Println("Failed to parse go.mod:", err)

		return ""
	}

	return modFile.Module.Mod.Path
}
