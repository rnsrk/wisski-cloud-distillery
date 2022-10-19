package drush

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errCronFailed = exit.Error{
	Message:  "Failed to run cron script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (drush *Drush) Cron(io stream.IOStream) error {
	code, err := drush.Barrel.Shell(io, "/runtime/cron.sh")
	if err != nil {
		io.EPrintln(err)
	}
	if code != 0 {
		// keep going, because we want to run as many crons as possible
		err = errCronFailed.WithMessageF(drush.Slug, code)
		io.EPrintln(err)
	}

	return nil
}

func (drush *Drush) LastCron(server *phpx.Server) (t time.Time, err error) {
	var timestamp int64
	err = drush.PHP.EvalCode(server, &timestamp, `$val = \Drupal::state()->get('system.cron_last'); return $val; `)
	if err != nil {
		return
	}
	return time.Unix(timestamp, 0), nil
}

type LastCronFetcher struct {
	ingredient.Base

	Drush *Drush
}

func (lbr *LastCronFetcher) Fetch(flags ingredient.FetchFlags, info *ingredient.Information) (err error) {
	if flags.Quick {
		return
	}

	info.LastRebuild, _ = lbr.Drush.LastCron(flags.Server)
	return
}
