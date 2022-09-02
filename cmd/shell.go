package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Shell is the 'shell' command
var Shell wisski_distillery.Command = shell{}

type shell struct {
	Positionals struct {
		Slug string   `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show run shell in"`
		Args []string `positional-arg-name:"ARGS" description:"arguments to pass to the shell"`
	} `positional-args:"true"`
}

func (shell) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsConfig: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "shell",
		Description: "Open a shell in the provided instance",
	}
}

var errShell = exit.Error{
	Message:  "Unable to start shell: %s",
	ExitCode: exit.ExitGeneric,
}

func (sh shell) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instance(sh.Positionals.Slug)
	if err != nil {
		return err
	}

	code, err := instance.Shell(context.IOStream, sh.Positionals.Args...)
	if err != nil {
		return errShell.WithMessageF(err)
	}
	if code != 0 {
		return exit.Error{
			ExitCode: exit.ExitCode(uint8(code)),
			Message:  fmt.Sprintf("Exit code %d", code),
		}
	}
	return nil
}
