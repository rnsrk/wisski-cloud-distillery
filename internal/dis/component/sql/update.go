package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

// Shell runs a mysql shell with the provided databases.
//
// NOTE(twiesing): This command should not be used to connect to the database or execute queries except in known situations.
func (sql *SQL) Shell(ctx context.Context, io stream.IOStream, argv ...string) (int, error) {
	return sql.Stack(sql.Environment).Exec(ctx, io, "sql", "mysql", argv...)
}

// unsafeWaitShell waits for a connection via the database shell to succeed
func (sql *SQL) unsafeWaitShell(ctx context.Context) error {
	n := stream.FromNil()
	return timex.TickUntilFunc(func(time.Time) bool {
		code, err := sql.Shell(ctx, n, "-e", "select 1;")
		return err == nil && code == 0
	}, ctx, sql.PollInterval)
}

// unsafeQuery shell executes a raw database query.
func (sql *SQL) unsafeQueryShell(ctx context.Context, query string) bool {
	code, err := sql.Shell(ctx, stream.FromNil(), "-e", query)
	return err == nil && code == 0
}

var errSQLUnableToCreateUser = errors.New("unable to create administrative user")
var errSQLUnsafeDatabaseName = errors.New("distillery database has an unsafe name")
var errSQLUnableToMigrate = exit.Error{
	Message:  "unable to migrate %s table: %s",
	ExitCode: exit.ExitGeneric,
}

// Update initializes or updates the SQL database.
func (sql *SQL) Update(ctx context.Context, io stream.IOStream) error {

	// unsafely create the admin user!
	{
		if err := sql.unsafeWaitShell(ctx); err != nil {
			return err
		}
		logging.LogMessage(io, "Creating administrative user")
		{
			username := sql.Config.MysqlAdminUser
			password := sql.Config.MysqlAdminPassword
			if err := sql.CreateSuperuser(ctx, username, password, true); err != nil {
				return errSQLUnableToCreateUser
			}
		}
	}

	// create the admin user
	logging.LogMessage(io, "Creating sql database")
	{
		if !sqle.IsSafeDatabaseLiteral(sql.Config.DistilleryDatabase) {
			return errSQLUnsafeDatabaseName
		}
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", sql.Config.DistilleryDatabase)
		if err := sql.Exec(createDBSQL); err != nil {
			return err
		}
	}

	// wait for the database to come up
	logging.LogMessage(io, "Waiting for database update to be complete")
	sql.WaitQueryTable(ctx)

	tables := []struct {
		name  string
		model any
		table string
	}{
		{
			"instance",
			&models.Instance{},
			models.InstanceTable,
		},
		{
			"metadata",
			&models.Metadatum{},
			models.MetadataTable,
		},
		{
			"snapshot",
			&models.Export{},
			models.ExportTable,
		},
		{
			"lock",
			&models.Lock{},
			models.LockTable,
		},
	}

	// migrate all of the tables!
	return logging.LogOperation(func() error {
		for _, table := range tables {
			logging.LogMessage(io, "migrating %q table", table.name)
			db, err := sql.QueryTable(ctx, false, table.table)
			if err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.name, "unable to access table")
			}

			if err := db.AutoMigrate(table.model); err != nil {
				return errSQLUnableToMigrate.WithMessageF(table.name, err)
			}
		}
		return nil
	}, io, "migrating database tables")
}
