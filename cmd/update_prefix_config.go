package cmd

import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/status"
)

// Cron is the 'cron' command
var UpdatePrefixConfig wisski_distillery.Command = updateprefixconfig{}

type updateprefixconfig struct {
	Parallel int `short:"p" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`
}

func (updateprefixconfig) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "update_prefix_config",
		Description: "updates the prefix configuration",
	}
}

var errPrefixUpdateFailed = exit.Error{
	Message:  "failed to update the prefix configuration",
	ExitCode: exit.ExitGeneric,
}

func (upc updateprefixconfig) Run(context wisski_distillery.Context) (err error) {
	defer errPrefixUpdateFailed.DeferWrap(&err)

	dis := context.Environment

	wissKIs, err := dis.Instances().All(context.Context)
	if err != nil {
		return err
	}

	return status.WriterGroup(context.Stderr, upc.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		io.WriteString(writer, "reading prefixes")
		err := instance.Prefixes().Update(context.Context)
		if err != nil {
			return err
		}
		return nil
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("update_prefix %q", item.Slug)
	}))
}
