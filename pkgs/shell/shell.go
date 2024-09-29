//go:generate mockgen -package mock_shell -destination=mock_shell/mock_shell.go . CmdInterface,Interface

package shell

import (
	"os/exec"
)

// CmdInterface represents an interface for exec.Cmd.
type CmdInterface interface {
	Start() error
	Run() error
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
}

// Command implements the CmdInterface using exec.Cmd.
type Command struct {
	cmd *exec.Cmd
}

// nolint:wrapcheck
func (c *Command) Start() error {
	return c.cmd.Start()
}

// nolint:wrapcheck
func (c *Command) Run() error {
	return c.cmd.Run()
}

// nolint:wrapcheck
func (c *Command) Output() ([]byte, error) {
	return c.cmd.Output()
}

// nolint:wrapcheck
func (c *Command) CombinedOutput() ([]byte, error) {
	return c.cmd.CombinedOutput()
}

type Shell struct{}

// nolint:ireturn
// NewCommand creates a new Command instance with the given name and arguments.
func (s *Shell) NewCommand(name string, args ...string) CmdInterface {
	cmd := exec.Command(name, args...)

	return &Command{cmd: cmd}
}

type Interface interface {
	NewCommand(name string, args ...string) CmdInterface
}
