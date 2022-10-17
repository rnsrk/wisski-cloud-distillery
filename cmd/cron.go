package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

// Cron is the 'cron' command
var Cron wisski_distillery.Command = cron{}

type cron struct {
	Parallel int `short:"p" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`

	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run cron in"`
	} `positional-args:"true"`
}

func (cron) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "cron",
		Description: "Runs the cron script for several instances",
	}
}

func (cr cron) Run(context wisski_distillery.Context) error {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(cr.Positionals.Slug...)
	if err != nil {
		return err
	}

	// and do the actual blind_update!
	return status.StreamGroup(context.IOStream, cr.Parallel, func(instance *wisski.WissKI, io stream.IOStream) error {
		return instance.Cron(io)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("cron %q", item.Slug)
	}))
}
