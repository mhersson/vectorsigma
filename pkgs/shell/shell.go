//go:generate mockgen -package mock_shell -destination=mock_shell/mock_shell.go . CmdRunner,Interface

package shell

import (
	"os/exec"
)

// CmdRunner represents an interface for exec.Cmd.
type CmdRunner interface {
	Run() error
}

// Command implements the CmdInterface using exec.Cmd.
type Command struct {
	cmd *exec.Cmd
}

// nolint:wrapcheck
func (c *Command) Run() error {
	return c.cmd.Run()
}

type Interface interface {
	NewCommand(name string, args ...string) CmdRunner
}

type Shell struct{}

// nolint:ireturn
// NewCommand creates a new Command instance with the given name and arguments.
func (s *Shell) NewCommand(name string, args ...string) CmdRunner {
	cmd := exec.Command(name, args...)

	return &Command{cmd: cmd}
}
