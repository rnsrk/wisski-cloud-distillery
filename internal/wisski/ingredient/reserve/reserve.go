package reserve

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
)

// Reserve implements reserving a WissKI Instance
// TODO: This should be integrated into the bookkeeping table.
type Reserve struct {
	ingredient.Base
}

//go:embed all:reserve reserve.env
var reserveResources embed.FS

// Stack returns a stack representing the reserve instance
func (reserve *Reserve) Stack() component.StackWithResources {
	return component.StackWithResources{
		Stack: component.Stack{
			Dir: reserve.FilesystemBase,
			Env: reserve.Malt.Environment,
		},

		Resources:   reserveResources,
		ContextPath: filepath.Join("reserve"),
		EnvPath:     filepath.Join("reserve.env"),

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": reserve.Malt.Config.DockerNetworkName,

			"SLUG":          reserve.Slug,
			"VIRTUAL_HOST":  reserve.Domain(),
			"HTTPS_ENABLED": reserve.Malt.Config.HTTPSEnabledEnv(),
		},
	}
}
