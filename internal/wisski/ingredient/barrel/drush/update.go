package drush

import (
	"context"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

var errBlindUpdateFailed = exit.Error{
	Message:  "Failed to run blind update script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

// Update performs a blind drush update
func (drush *Drush) Update(ctx context.Context, io stream.IOStream) error {
	code, err := drush.Barrel.Shell(ctx, io, "/runtime/blind_update.sh")
	if err != nil {
		return errBlindUpdateFailed.WithMessageF(drush.Slug, environment.ExecCommandError)
	}
	if code != 0 {
		return errBlindUpdateFailed.WithMessageF(drush.Slug, code)
	}

	return drush.setLastUpdate(ctx)
}

const lastUpdate = mstore.For[int64]("lastUpdate")

func (drush *Drush) LastUpdate(ctx context.Context) (t time.Time, err error) {
	epoch, err := lastUpdate.Get(ctx, drush.MStore)
	if err == meta.ErrMetadatumNotSet {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	// and turn it into time!
	return time.Unix(epoch, 0), nil
}

func (drush *Drush) setLastUpdate(ctx context.Context) error {
	return lastUpdate.Set(ctx, drush.MStore, time.Now().Unix())
}

type LastUpdateFetcher struct {
	ingredient.Base

	Drush *Drush
}

var (
	_ ingredient.WissKIFetcher = (*LastUpdateFetcher)(nil)
)

func (lbr *LastUpdateFetcher) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	info.LastUpdate, err = lbr.Drush.LastUpdate(flags.Context)
	return
}
