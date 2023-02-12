package instances

import (
	"context"
	"embed"
	"io"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/tkw1536/goprogram/exit"
)

var errBootstrapFailedRuntime = exit.Error{
	Message:  "failed to update runtime",
	ExitCode: exit.ExitGeneric,
}

// Runtime contains runtime resources to be installed into any instance
//
//go:embed all:runtime
var runtimeResources embed.FS

// Update installs or updates runtime components needed by this component.
func (instances *Instances) Update(ctx context.Context, progress io.Writer) error {
	err := unpack.InstallDir(instances.Still.Environment, instances.Config.Paths.RuntimeDir(), "runtime", runtimeResources, func(dst, src string) {
		logging.ProgressF(progress, ctx, "[copy]  %s\n", dst)
	})
	if err != nil {
		return errBootstrapFailedRuntime.Wrap(err)
	}
	return nil
}
