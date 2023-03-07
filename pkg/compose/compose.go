package compose

import (
	"path/filepath"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

// ComposeProject represents a compose project
type Project = *types.Project

// Open loads a docker compose project from the given path.
// The returned name indicates the name, as would be found by the 'docker compose' executable.
// If the project could not be found, an appropriate error is returned.
//
// NOTE: This intentionally omits using any kind of api for docker compose.
// This saves a *a lot* of dependencies.
func Open(path string) (project Project, err error) {
	ppath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	opts, err := cli.NewProjectOptions(
		/* configs = */ nil,
		cli.WithWorkingDirectory(ppath),
		cli.WithDefaultConfigPath,
		cli.WithName(loader.NormalizeProjectName(filepath.Base(ppath))),
		cli.WithDotEnv,
	)
	if err != nil {
		return nil, err
	}

	proj, err := cli.ProjectFromOptions(opts)
	if err != nil {
		return nil, err
	}

	return proj, nil
}
