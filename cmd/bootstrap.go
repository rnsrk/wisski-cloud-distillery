package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	cfg "github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
	"github.com/FAU-CDI/wisski-distillery/internal/hostname"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/password"
	"github.com/tkw1536/goprogram/exit"
)

// Bootstrap is the 'bootstrap' command
var Bootstrap wisski_distillery.Command = bootstrap{}

type bootstrap struct {
	Directory string `short:"r" long:"root-directory" description:"path to the root deployment directory" default:"/var/www/deploy"`
	Hostname  string `short:"h" long:"hostname" description:"default hostname of the distillery (default: system hostname)"`
}

func (bootstrap) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: env.Requirements{
			NeedsDistillery: false,
		},
		Command:     "bootstrap",
		Description: "Bootstraps the installation of a Distillery System",
	}
}

var errBootstrapDifferent = exit.Error{
	Message:  "refusing to bootstrap: base directory is already set to %s.",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapFailedToCreateDirectory = exit.Error{
	Message:  "failed to create directory %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapFailedToSaveDirectory = exit.Error{
	Message:  "failed to register base directory: %s",
	ExitCode: exit.ExitGeneric,
}

var errBoostrapFailedToCopyExe = exit.Error{
	Message:  "failed to copy wdcli executable: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapWriteConfig = exit.Error{
	Message:  "failed to write configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapOpenConfig = exit.Error{
	Message:  "failed to open configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapCreateFile = exit.Error{
	Message:  "failed to touch configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

func (bs bootstrap) Run(context wisski_distillery.Context) error {
	root := bs.Directory

	// check that we didn't get a different base directory
	{
		got, err := env.ReadBaseDirectory()
		if err == nil && got != "" && got != root {
			return errBootstrapDifferent.WithMessageF(got)
		}
	}

	{
		logging.LogMessage(context.IOStream, "Creating root deployment directory")
		if err := os.MkdirAll(root, fs.ModeDir); err != nil {
			return errBootstrapFailedToCreateDirectory.WithMessageF(root)
		}
		if err := env.WriteBaseDirectory(root); err != nil {
			return errBootstrapFailedToSaveDirectory.WithMessageF(root)
		}
		context.Println(root)
	}

	// TODO: Read these from the command line?
	wdcliPath := filepath.Join(root, env.Executable)
	envPath := filepath.Join(root, env.ConfigFile)
	domain := bs.Hostname
	if domain == "" {
		domain = hostname.FQDN()
	}
	overridesPath := filepath.Join(root, "overrides.json")
	authorizedKeysFile := filepath.Join(root, "authorized_keys")

	{
		logging.LogMessage(context.IOStream, "Copying over wdcli executable")
		exe, err := os.Executable()
		if err != nil {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}

		err = fsx.CopyFile(wdcliPath, exe)
		if err != nil && err != fsx.ErrCopySameFile {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}
		context.Println(wdcliPath)
	}

	{
		if !fsx.IsFile(envPath) {
			if err := logging.LogOperation(func() error {
				password, err := password.Password(128)
				if err != nil {
					return errBootstrapWriteConfig.WithMessageF(err)
				}

				if err := distillery.InstallTemplate(envPath, filepath.Join("resources", "templates", "bootstrap", "env"), map[string]string{
					"DEPLOY_ROOT":          root,
					"DEFAULT_DOMAIN":       domain,
					"SELF_OVERRIDES_FILE":  overridesPath,
					"AUTHORIZED_KEYS_FILE": authorizedKeysFile,

					"GRAPHDB_ADMIN_USER":     "admin",
					"GRAPHDB_ADMIN_PASSWORD": password[:64],

					"MYSQL_ADMIN_USER":     "admin",
					"MYSQL_ADMIN_PASSWORD": password[64:],
				}); err != nil {
					return errBootstrapWriteConfig.WithMessageF(err)
				}

				return nil
			}, context.IOStream, "Installing configuration file"); err != nil {
				return err
			}

			if err := logging.LogOperation(func() error {

				context.Println(overridesPath)
				if err := distillery.InstallTemplate(overridesPath, filepath.Join("resources", "templates", "bootstrap", "overrides.json"), map[string]string{}); err != nil {
					return errBootstrapCreateFile.WithMessageF(err)
				}

				context.Println(authorizedKeysFile)
				if err := distillery.InstallTemplate(authorizedKeysFile, filepath.Join("resources", "templates", "bootstrap", "global_authorized_keys"), map[string]string{}); err != nil {
					return errBootstrapCreateFile.WithMessageF(err)
				}

				return nil
			}, context.IOStream, "Creating additional config files"); err != nil {
				return err
			}
		}

	}

	// re-read the configuration and print it!
	logging.LogMessage(context.IOStream, "Configuration is now complete")
	f, err := os.Open(envPath)
	if err != nil {
		return errBootstrapOpenConfig.WithMessageF(err)
	}
	defer f.Close()

	var config cfg.Config
	if err := config.Unmarshal(f); err != nil {
		return errBootstrapOpenConfig.WithMessageF(err)
	}
	context.Println(config)

	// Tell the user how to proceed
	logging.LogMessage(context.IOStream, "Bootstrap is complete")
	context.Printf("Adjust the configuration file at %s\n", envPath)
	context.Printf("Then grab a GraphDB zipped source file and run:\n")
	context.Printf("%s system_update /path/to/graphdb.zip\n", wdcliPath)

	return nil
}
