package config

import (
	"crypto/rand"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/pkglib/hostname"
	"github.com/tkw1536/pkglib/password"
)

// Template is used to generate a configuration file.
type Template struct {
	RootPath      string
	DefaultDomain string

	TSAdminUser     string
	TSAdminPassword string

	SQLAdminUsername string
	SQLAdminPassword string

	DockerNetworkName string
	SessionSecret     string
}

// SetDefaults sets defaults on the template
func (tpl *Template) SetDefaults(env environment.Environment) (err error) {
	if tpl.RootPath == "" {
		tpl.RootPath = bootstrap.BaseDirectoryDefault
	}

	if tpl.DefaultDomain == "" {
		tpl.DefaultDomain = hostname.FQDN() // TODO: Make this environment specific
	}

	if tpl.TSAdminUser == "" {
		tpl.TSAdminUser = "admin"
	}

	if tpl.TSAdminPassword == "" {
		tpl.TSAdminPassword, err = password.Generate(rand.Reader, 64, passwordx.Safe)
		if err != nil {
			return err
		}
	}

	if tpl.SQLAdminUsername == "" {
		tpl.SQLAdminUsername = "admin"
	}

	if tpl.SQLAdminPassword == "" {
		tpl.SQLAdminPassword, err = password.Generate(rand.Reader, 64, passwordx.Safe)
		if err != nil {
			return err
		}
	}

	if tpl.DockerNetworkName == "" {
		tpl.DockerNetworkName, err = password.Generate(rand.Reader, 10, passwordx.Identifier)
		if err != nil {
			return err
		}
		tpl.DockerNetworkName = `distillery-` + tpl.DockerNetworkName
	}

	if tpl.SessionSecret == "" {
		tpl.SessionSecret, err = password.Generate(rand.Reader, 100, passwordx.Printable)
		if err != nil {
			return err
		}
	}

	return nil
}

// Generate generates a configuration file for this configuration
func (tpl Template) Generate() Config {
	return Config{
		Paths: PathsConfig{
			Root:           tpl.RootPath,
			OverridesJSON:  filepath.Join(tpl.RootPath, bootstrap.OverridesJSON),
			ResolverBlocks: filepath.Join(tpl.RootPath, bootstrap.ResolverBlockedTXT),
		},
		HTTP: HTTPConfig{
			PrimaryDomain: tpl.DefaultDomain,
			ExtraDomains:  []string{},
		},
		Docker: DockerConfig{
			tpl.DockerNetworkName,
		},
		SQL: SQLConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.SQLAdminUsername,
				AdminPassword: tpl.SQLAdminPassword,

				UserPrefix: "mysql-factory-",
				DataPrefix: "mysql-factory-",
			},

			Database: "distillery",
		},
		TS: TSConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.TSAdminUser,
				AdminPassword: tpl.TSAdminPassword,

				UserPrefix: "graphdb-factory-",
				DataPrefix: "graphdb-factory-",
			},
		},
		MaxBackupAge:   30 * 24 * time.Hour, // 1 month
		PasswordLength: 64,

		PublicSSHPort: 2222,

		SessionSecret: tpl.SessionSecret,
		CronInterval:  10 * time.Minute,
	}
}
