package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
	"github.com/FAU-CDI/wisski-distillery/pkg/smartp"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

// BlindUpdate is the 'blind_update' command
var BlindUpdate wisski_distillery.Command = blindUpdate{}

type blindUpdate struct {
	Parallel    int  `short:"p" long:"parallel" description:"run on (at most) this many instances concurrently. 0 for no limit." default:"1"`
	Force       bool `short:"f" long:"force" description:"force running blind-update even if AutoBlindUpdate is set to false"`
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run blind-update in"`
	} `positional-args:"true"`
}

func (blindUpdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "blind_update",
		Description: "Runs the blind update in the provided instances",
	}
}

var errBlindUpdateFailed = exit.Error{
	Message:  "Failed to run blind update script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (bu blindUpdate) Run(context wisski_distillery.Context) error {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(bu.Positionals.Slug...)
	if err != nil {
		return err
	}
	if !bu.Force {
		wissKIs = slicesx.Filter(wissKIs, func(instance instances.WissKI) bool {
			return bool(instance.AutoBlindUpdateEnabled)
		})
	}

	// and do the actual blind_update!
	return smartp.Run(context.IOStream, bu.Parallel, func(instance instances.WissKI, io stream.IOStream) error {
		code, err := instance.Shell(io, "/runtime/blind_update.sh")
		if err != nil {
			return errBlindUpdateFailed.WithMessageF(instance.Slug, environment.ExecCommandError)
		}
		if code != 0 {
			return errBlindUpdateFailed.WithMessageF(instance.Slug, code)
		}
		return nil
	}, wissKIs, smartp.SmartMessage(func(item instances.WissKI) string {
		return fmt.Sprintf("blind_update %q", item.Slug)
	}))
}
