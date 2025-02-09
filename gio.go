package gio

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Command embeds an exec.Cmd and also includes buffers for std{in,err,out}
// These will automatically attached when they're instantiated.
type Command struct {
	*exec.Cmd
	BioStdin, BioStdout, BioStderr *bytes.Buffer // Bio from Plan9's buffer io (Bio)
}

// Satisfies the interface
func (c *Command) String() string {
	return strings.Join(c.Args, " ")
}

// Constructor: i.e. instantiates a new pointer to a Command struct and attaches the command's
// std{in,out,err} file descriptors
func New(name string, arg ..string) *Command {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command(name, arg...)
	cmd.Stdin := stdin
	cmd.Stdout := stdout
	cmd.Stderr:= stderr

	return &Command{
		Cmd:		cmd,
		BioStdin:	stdin,
		BioStdout:	stdout,
		BioStderr:	stderr,
	}
}

// Run calls `*exec.Cmd.Run()` as a wrapper to Command and returns
// an error. If `Run()` fails, Run() will return the error also by
// checking stderr's buffer.
//
// If stderr's buffer is empty, an error was returned with its contents.
func (c *Command) Run() error {
	if err := c.Start(); err != nil { return fmt.Errorf("[ERROR]: cmd %s failed with: %v", c, err) }
	if err := c.Wait(); err != nil { return err }
	return nil
}

// Wait calls `*exec.Cmd.Wait()` as a wrapper to Command and handles
// errors per Run().
func (c *Command) Wait() error {
	if err := c.Cmd.Wait(); err != nil {
		if c.BioStderr.Len() > 0 {
			return fmt.Errorf("[ERROR]: cmd %s: failes with %s.\n\n here:%s",
				c, err, c.BioStderr.String())
		}
		return fmt.Errorf("[ERROR]: running %s failed with: %s.", c, err)
	}
	return nil
}
