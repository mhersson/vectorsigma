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
package cmd

import (
	"fmt"
	"os"

	"github.com/mhersson/vectorsigma/pkgs/generator"

	"github.com/spf13/cobra"
)

const (
	initFlag    = "init"
	moduleFlag  = "module"
	inputFlag   = "input"
	outputFlag  = "output"
	packageFlag = "package"
)

// GenerateCmd represents the generate command.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a FSM",
	Long:  `Generate a FSM base on a UML input file.`,
	Run: func(cmd *cobra.Command, _ []string) {
		init := cmd.Flag(initFlag).Changed
		input, err := cmd.Flags().GetString(inputFlag)
		generator.CheckError(err, "Error reading input flag")
		output, err := cmd.Flags().GetString(outputFlag)
		generator.CheckError(err, "Error reading output flag")
		pkg, err := cmd.Flags().GetString(packageFlag)
		generator.CheckError(err, "Error reading package flag")
		module, err := cmd.Flags().GetString(moduleFlag)
		generator.CheckError(err, "Error reading module flag")

		if input == "" {
			fmt.Println(generator.ErrorPrefix + " missing required flag input")
			_ = cmd.Usage()
			os.Exit(1)
		}

		generator.Run(init, module, pkg, input, output)
	},
}

func init() {
	rootCmd.AddCommand(GenerateCmd)

	GenerateCmd.Flags().Bool(initFlag, false, "Initialize new go module")
	GenerateCmd.Flags().StringP(moduleFlag, "m", "", "Name of new go module (defaults to current directory name)")
	GenerateCmd.Flags().StringP(inputFlag, "i", "", "The UML input file")
	GenerateCmd.Flags().StringP(outputFlag, "o", "",
		"The output path of the generated FSM (defaults to current directory)")
	GenerateCmd.Flags().StringP(packageFlag, "p", "fsm", "The package name of the generated FSM")
}
