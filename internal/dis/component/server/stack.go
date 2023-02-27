package server

import (
	"context"
	"embed"
	"io"
	"path/filepath"
	"syscall"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

func (control Server) Path() string {
	return filepath.Join(control.Still.Config.Paths.Root, "core", "dis")
}

//go:embed all:server server.env
var resources embed.FS

func (server *Server) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(server, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "server",
		EnvPath:     "server.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": server.Config.Docker.Network,
			"HOST_RULE":           server.Config.HTTP.DefaultHostRule(),
			"HTTPS_ENABLED":       server.Config.HTTP.HTTPSEnabledEnv(),

			"CONFIG_PATH": server.Config.ConfigPath,
			"DEPLOY_ROOT": server.Config.Paths.Root,

			"SELF_OVERRIDES_FILE":      server.Config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": server.Config.Paths.ResolverBlocks,

			"CUSTOM_ASSETS_PATH": server.Dependencies.Templating.CustomAssetsPath(),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

// Trigger triggers the active cron run to immediatly invoke cron.
func (server *Server) Trigger(ctx context.Context, env environment.Environment) error {
	return server.Stack(env).Kill(ctx, io.Discard, "control", syscall.SIGHUP)
}

func (server *Server) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: server.Config.Paths.CurrentExecutable(server.Environment), // TODO: Does this make sense?
	}
}