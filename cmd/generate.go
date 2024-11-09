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

var inputParams *generator.InputParams

// GenerateCmd represents the generate command.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a FSM",
	Long:  `Generate a FSM base on a UML input file.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return generator.Run(inputParams)
	},
}

func init() {
	inputParams = &generator.InputParams{}

	rootCmd.AddCommand(GenerateCmd)

	GenerateCmd.Flags().BoolVar(&inputParams.Init, initFlag, false, "Initialize new go module")
	GenerateCmd.Flags().StringVarP(&inputParams.Module, moduleFlag, "m", "",
		"Name of new go module (default current directory name)")
	GenerateCmd.Flags().StringVarP(&inputParams.Input, inputFlag, "i", "", "The UML input file")
	_ = GenerateCmd.MarkFlagRequired(inputFlag)
	GenerateCmd.Flags().StringVarP(&inputParams.Output, outputFlag, "o", "",
		"The output path of the generated FSM (default current working directory)")
	GenerateCmd.Flags().StringVarP(&inputParams.Package, packageFlag, "p", "fsm", "The package name of the generated FSM")
}
