package triplestore

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"
)

type Triplestore struct {
	component.Base

	BaseURL string // upstream server url

	PollInterval time.Duration // duration to wait for during wait
}

var (
	_ component.Backupable    = (*Triplestore)(nil)
	_ component.Snapshotable  = (*Triplestore)(nil)
	_ component.Installable   = (*Triplestore)(nil)
	_ component.Provisionable = (*Triplestore)(nil)
)

func (ts *Triplestore) Path() string {
	return filepath.Join(ts.Still.Config.Paths.Root, "core", "triplestore")
}

func (Triplestore) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:triplestore
var resources embed.FS

func (ts *Triplestore) Stack() component.StackWithResources {
	return component.MakeStack(ts, component.StackWithResources{
		Resources:   resources,
		ContextPath: "triplestore",

		CopyContextFiles: []string{"graphdb.zip"}, // TODO: Move into constant?

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ts.Config.Docker.Network(),
			"HOST_RULE":           ts.Config.HTTP.HostRule(config.TriplestoreDomain.Domain()),
			"HTTPS_ENABLED":       ts.Config.HTTP.HTTPSEnabledEnv(),
		},

		ComposerYML: func(root *yaml.Node) (*yaml.Node, error) {
			// ts is exposed => everything is fine
			if ts.Config.HTTP.TS.Set && ts.Config.HTTP.TS.Value {
				return root, nil
			}

			// not exposed => remove the appropriate labels
			if err := yamlx.ReplaceWith(root, []string{
				"eu.wiss-ki.barrel.distillery=${DOCKER_NETWORK_NAME}",
			}, "services", "triplestore", "labels"); err != nil {
				return nil, err
			}

			return root, nil
		},

		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
			filepath.Join("data", "import"),
		},
	})
}
