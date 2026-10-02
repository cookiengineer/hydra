package lib

import "bytes"
import "os"
import "os/exec"
import "strings"

type Command struct {
	Name string
	Args []string
	Env  []string
	Dir  string
}

type CommandResult struct {
	Command  string
	Stdout   string
	Stderr   string
	ExitCode int
}

type Executor interface {
	Run(command Command) CommandResult
}

type OSExecutor struct{}

func (executor OSExecutor) Run(command Command) CommandResult {

	cmd := exec.Command(command.Name, command.Args...)

	cmd.Env = append(os.Environ(), command.Env...)
	cmd.Dir = command.Dir

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	code := 0

	if err != nil {

		if exit_err, ok := err.(*exec.ExitError); ok {
			code = exit_err.ExitCode()
		} else {
			code = -1
		}

	}

	return CommandResult{
		Command:  command.Name + " " + strings.Join(command.Args, " "),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: code,
	}

}
