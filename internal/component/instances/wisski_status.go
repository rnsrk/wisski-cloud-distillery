package instances

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

// WissKIInfo represents information about this WissKI Instance.
type WissKIInfo struct {
	Time time.Time // Time this info was built

	// Generic Information
	Slug string // slug
	URL  string // complete URL, including http(s)

	// Information about the running instance
	Running     bool
	LastRebuild time.Time

	// List of backups made
	Snapshots []models.Export

	// WissKI content information
	NoPrefixes   bool              // TODO: Move this into the database
	Prefixes     []string          // list of prefixes
	Pathbuilders map[string]string // all the pathbuilders
}

// Info generate a
func (wisski *WissKI) Info(quick bool) (info WissKIInfo, err error) {
	// TODO: Cache this, and run it with every cron!

	info.Time = time.Now().UTC()

	// static properties
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()

	// dynamic properties
	var group errgroup.Group

	// quick check if this wisski is running
	group.Go(func() (err error) {
		info.Running, err = wisski.Running()
		return
	})

	// slower checks for extra properties.
	// these might execute php code or require additional database queries.
	if !quick {
		group.Go(func() (err error) {
			info.LastRebuild, _ = wisski.LastRebuild()
			return nil
		})
		group.Go(func() error {
			info.Pathbuilders, _ = wisski.AllPathbuilders()
			return nil
		})
		group.Go(func() (err error) {
			info.Prefixes, _ = wisski.Prefixes()
			info.NoPrefixes = wisski.NoPrefix()
			return nil
		})
		group.Go(func() error {
			info.Snapshots, _ = wisski.Snapshots()
			return nil
		})
	}

	err = group.Wait()
	return
}

// Running checks if this WissKI is currently running.
func (wisski *WissKI) Running() (bool, error) {
	ps, err := wisski.Barrel().Ps(stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}
