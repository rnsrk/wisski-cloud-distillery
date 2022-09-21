package sql

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/sqle"
)

var errProvisionInvalidDatabaseParams = errors.New("Provision: Invalid parameters")
var errProvisionInvalidGrant = errors.New("Provision: Grant failed")

// ProvisionInstance provisions sql-specific resource for the given instance
func (sql *SQL) Provision(instance models.Instance, domain string) error {
	return sql.CreateDatabase(instance.SqlDatabase, instance.SqlUsername, instance.SqlPassword)
}

// CreateDatabase creates a new database with the given name.
// It then generates a new user, with the name 'user' and the password 'password', that is then granted access to this database.
//
// Provision internally waits for the database to become available.
func (sql *SQL) CreateDatabase(name, user, password string) error {

	// NOTE(twiesing): We shouldn't use string concat to build sql queries.
	// But the driver doesn't support using query params for this particular query.
	// Apparently it's a "feature", see https://github.com/go-sql-driver/mysql/issues/398#issuecomment-169951763.

	// quick and dirty check to make sure that all the names won't sql inject.
	if !sqle.IsSafeDatabaseLiteral(name) || !sqle.IsSafeDatabaseSingleQuote(user) || !sqle.IsSafeDatabaseSingleQuote(password) {
		return errProvisionInvalidDatabaseParams
	}

	// We use the sql shell here, because not only can we not use query params, but the driver outright rejects queries.
	// Queries of the form "CREATE USER 'test'@'%' IDENTIFIED BY 'test'; FLUSH PRIVILEGES;" return error 1064 when using driver, but are fine with the shell.
	// This should be fixed eventually, but I have no idea how.

	if err := sql.unsafeWaitShell(); err != nil {
		return err
	}

	query := "CREATE DATABASE `" + name + "`;" +
		"CREATE USER '" + user + "'@'%' IDENTIFIED BY '" + password + "';" +
		"GRANT ALL PRIVILEGES ON `" + name + "`.* TO `" + user + "`@`%`; FLUSH PRIVILEGES;"
	if !sql.unsafeQueryShell(query) {
		return errProvisionInvalidGrant
	}

	return nil
}

var errCreateSuperuserGrant = errors.New("CreateSuperUser: Grant failed")

// CreateSuperuser createsa new user, with the name 'user' and the password 'password'.
// It then grants this user superuser status in the database.
//
// CreateSuperuser internally waits for the database to become available.
func (sql *SQL) CreateSuperuser(user, password string, allowExisting bool) error {
	// NOTE(twiesing): This function unsafely uses the shell directly to create a superuser.
	// This is for two reasons:
	// (1) this is used during bootstraping
	// (2) The underlying driver doesn't support "GRANT ALL PRIVILEGES"
	// See also [sql.Provision].

	if !sqle.IsSafeDatabaseSingleQuote(user) || !sqle.IsSafeDatabaseSingleQuote(password) {
		return errProvisionInvalidDatabaseParams
	}

	if err := sql.unsafeWaitShell(); err != nil {
		return err
	}

	var IfNotExists string
	if allowExisting {
		IfNotExists = "IF NOT EXISTS"
	}

	query := "CREATE USER " + IfNotExists + " '" + user + "'@'%' IDENTIFIED BY '" + password + "';" +
		"GRANT ALL PRIVILEGES ON *.* TO '" + user + "'@'%' WITH GRANT OPTION; FLUSH PRIVILEGES;"
	if !sql.unsafeQueryShell(query) {
		return errCreateSuperuserGrant
	}

	return nil
}

var errPurgeUser = errors.New("PurgeUser: Failed to drop user")

// SQLPurgeUser deletes the specified user from the database
func (sql *SQL) PurgeUser(user string) error {
	if !sqle.IsSafeDatabaseSingleQuote(user) {
		return errPurgeUser
	}

	query := "DROP USER IF EXISTS '" + user + "'@'%';" +
		"FLUSH PRIVILEGES;"
	if !sql.unsafeQueryShell(query) {
		return errPurgeUser
	}

	return nil
}

var errSQLPurgeDB = errors.New("unable to drop database: unsafe database name")

// SQLPurgeDatabase deletes the specified db from the database
func (sql *SQL) PurgeDatabase(db string) error {
	if !sqle.IsSafeDatabaseLiteral(db) {
		return errSQLPurgeDB
	}
	return sql.Exec("DROP DATABASE IF EXISTS `" + db + "`")
}