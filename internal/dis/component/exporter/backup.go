package exporter

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/status"
	"golang.org/x/exp/slices"
)

// Backup describes a backup
type Backup struct {
	Description BackupDescription

	// Start and End Time of the backup
	StartTime time.Time
	EndTime   time.Time

	// various error states, which are ignored when creating the snapshot
	ErrPanic interface{}

	// errors for the various components
	ComponentErrors map[string]error

	// TODO: Make this proper
	ConfigFileErr error

	// Snapshots containing instances
	InstanceListErr   error
	InstanceSnapshots []Snapshot

	// List of files included
	WithManifest
}

// BackupDescription provides a description for a backup
type BackupDescription struct {
	Dest string // Destination path

	ConcurrentSnapshots int // maximum number of concurrent snapshots
}

// New create a new Backup
func (exporter *Exporter) NewBackup(ctx context.Context, progress io.Writer, description BackupDescription) (backup Backup) {
	backup.Description = description

	// catch anything critical that happened during the snapshot
	defer func() {
		backup.ErrPanic = recover()
	}()

	// do the create keeping track of time!
	logging.LogOperation(func() error {
		backup.StartTime = time.Now().UTC()
		backup.run(ctx, progress, exporter)
		backup.EndTime = time.Now().UTC()

		return nil
	}, progress, ctx, "Writing backup files")

	return
}

func (backup *Backup) run(ctx context.Context, progress io.Writer, exporter *Exporter) {
	// create a manifest
	manifest, done := backup.handleManifest(backup.Description.Dest)
	defer done()

	// create a new status display
	backups := exporter.Dependencies.Backupable
	backup.ComponentErrors = make(map[string]error, len(backups))

	// Component backup tasks
	logging.LogOperation(func() error {
		st := status.NewWithCompat(progress, 0)
		st.Start()
		defer st.Stop()

		errors := status.Group[component.Backupable, error]{
			PrefixString: func(item component.Backupable, index int) string {
				return fmt.Sprintf("[backup %q]: ", item.Name())
			},
			PrefixAlign: true,

			Handler: func(bc component.Backupable, index int, writer io.Writer) error {
				return bc.Backup(
					component.NewStagingContext(
						ctx,
						exporter.Environment,
						writer,
						filepath.Join(backup.Description.Dest, bc.BackupName()),
						manifest,
					),
				)
			},

			ResultString: status.DefaultErrorString[component.Backupable],
		}.Use(st, backups)

		for i, bc := range backups {
			backup.ComponentErrors[bc.Name()] = errors[i]
		}

		return nil
	}, progress, ctx, "Backing up core components")

	// backup instances
	logging.LogOperation(func() error {
		st := status.NewWithCompat(progress, 0)
		st.Start()
		defer st.Stop()

		instancesBackupDir := filepath.Join(backup.Description.Dest, "instances")
		if err := fsx.Mkdir(instancesBackupDir, fsx.DefaultDirPerm); err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// list all instances
		wissKIs, err := exporter.Dependencies.Instances.All(ctx)
		if err != nil {
			backup.InstanceListErr = err
			return nil
		}

		// make a backup of the snapshots
		backup.InstanceSnapshots = status.Group[*wisski.WissKI, Snapshot]{
			PrefixString: func(item *wisski.WissKI, index int) string {
				return fmt.Sprintf("[snapshot %q]: ", item.Slug)
			},
			PrefixAlign: true,

			Handler: func(instance *wisski.WissKI, index int, writer io.Writer) Snapshot {
				dir := filepath.Join(instancesBackupDir, instance.Slug)
				if err := fsx.Mkdir(dir, fsx.DefaultDirPerm); err != nil {
					return Snapshot{
						ErrPanic: err,
					}
				}

				manifest <- dir

				return exporter.NewSnapshot(ctx, instance, writer, SnapshotDescription{
					Dest: dir,
				})
			},
			ResultString: func(res Snapshot, item *wisski.WissKI, index int) string {
				return "done"
			},
			WaitString:   status.DefaultWaitString[*wisski.WissKI],
			HandlerLimit: backup.Description.ConcurrentSnapshots,
		}.Use(st, wissKIs)

		// sort the instances
		slices.SortFunc(backup.InstanceSnapshots, func(a, b Snapshot) bool {
			return a.Instance.Slug < b.Instance.Slug
		})

		return nil
	}, progress, ctx, "Creating instance snapshots")

}
