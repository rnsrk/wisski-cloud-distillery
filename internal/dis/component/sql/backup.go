package sql

import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/goprogram/stream"
)

var errSQLBackup = errors.New("`SQLBackup': mysqldump returned non-zero exit code")

func (*SQL) BackupName() string {
	return "sql.sql"
}

// Backup makes a backup of all SQL databases into the path dest.
func (sql *SQL) Backup(scontext component.StagingContext) error {
	return scontext.AddFile("", func(ctx context.Context, file io.Writer) error {
		code := sql.Stack().Exec(ctx, stream.NewIOStream(file, scontext.Progress(), nil, 0), "sql", "mysqldump", "--all-databases")()
		if code != 0 {
			return errSQLBackup
		}
		return nil
	})

}
