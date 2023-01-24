// Package stack implements a docker compose stack
package component

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/pools"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

// Stack represents a 'docker compose' stack living in a specific directory
//
// NOTE(twiesing): In the current implementation this requires a 'docker' executable on the system.
// This executable must be capable of the 'docker compose' command.
// In the future the idea is to replace this with a native docker compose client.
type Stack struct {
	Dir string // Directory this Stack is located in

	Env              environment.Environment
	DockerExecutable string // Path to the native docker executable to use
}

var errStackKill = errors.New("Stack.Kill: Kill returned non-zero exit code")

func (ds Stack) Kill(ctx context.Context, progress io.Writer, service string, signal os.Signal) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "kill", service, "-s", signal.String())()
	if code != 0 {
		return errStackKill
	}
	return nil
}

var errStackUpdatePull = errors.New("Stack.Update: Pull returned non-zero exit code")
var errStackUpdateBuild = errors.New("Stack.Update: Build returned non-zero exit code")

// Update pulls, builds, and then optionally starts this stack.
// This does not have a direct 'docker compose' shell equivalent.
//
// See also Up.
func (ds Stack) Update(ctx context.Context, progress io.Writer, start bool) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "pull")(); code != 0 {
		return errStackUpdatePull
	}

	if code := ds.compose(ctx, stream.NonInteractive(progress), "build", "--pull")(); code != 0 {
		return errStackUpdateBuild
	}

	if start {
		return ds.Up(ctx, progress)
	}
	return nil
}

var errStackUp = errors.New("Stack.Up: Up returned non-zero exit code")

// Up creates and starts the containers in this Stack.
// It is equivalent to 'docker compose up --force-recreate --remove-orphans --detach' on the shell.
func (ds Stack) Up(ctx context.Context, progress io.Writer) error {
	if code := ds.compose(ctx, stream.NonInteractive(progress), "up", "--force-recreate", "--remove-orphans", "--detach")(); code != 0 {
		return errStackUp
	}
	return nil
}

// Exec executes an executable in the provided running service.
// It is equivalent to 'docker compose exec $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Exec(ctx context.Context, io stream.IOStream, service, executable string, args ...string) func() int {
	compose := []string{"exec"}
	if io.StdinIsATerminal() {
		compose = append(compose, "-ti")
	}

	compose = append(compose, service)
	compose = append(compose, executable)
	compose = append(compose, args...)

	return ds.compose(ctx, io, compose...)
}

// Run runs a command in a running container with the given executable.
// It is equivalent to 'docker compose run [--rm] $service $executable $args...'.
//
// It returns the exit code of the process.
func (ds Stack) Run(ctx context.Context, io stream.IOStream, autoRemove bool, service, command string, args ...string) (int, error) {
	compose := []string{"run"}
	if autoRemove {
		compose = append(compose, "--rm")
	}
	if !io.StdinIsATerminal() {
		compose = append(compose, "-T")
	}
	compose = append(compose, service, command)
	compose = append(compose, args...)

	code := ds.compose(ctx, io, compose...)()
	return code, nil
}

var errStackRestart = errors.New("Stack.Restart: Restart returned non-zero exit code")

// Restart restarts all containers in this Stack.
// It is equivalent to 'docker compose restart' on the shell.
func (ds Stack) Restart(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "restart")()
	if code != 0 {
		return errStackRestart
	}
	return nil
}

var errStackPs = errors.New("Stack.Ps: Down returned non-zero exit code")

// Ps returns the ids of the containers currently running
func (ds Stack) Ps(ctx context.Context, progress io.Writer) ([]string, error) {
	// create a buffer
	buffer := pools.GetBuffer()
	defer pools.ReleaseBuffer(buffer)

	// read the ids from the command!
	code := ds.compose(ctx, stream.NewIOStream(buffer, progress, nil, 0), "ps", "-q")()
	if code != 0 {
		return nil, errStackPs
	}

	// scan each of the lines
	var results []string
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		if text := scanner.Text(); text != "" {
			results = append(results, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// return them!
	return results, nil
}

var errStackDown = errors.New("Stack.Down: Down returned non-zero exit code")

// Down stops and removes all containers in this Stack.
// It is equivalent to 'docker compose down -v' on the shell.
func (ds Stack) Down(ctx context.Context, progress io.Writer) error {
	code := ds.compose(ctx, stream.NonInteractive(progress), "down", "-v")()
	if code != 0 {
		return errStackDown
	}
	return nil
}

// compose executes a 'docker compose' command on this stack.
//
// NOTE(twiesing): Check if this can be replaced by an internal call to libcompose.
// But probably not.
func (ds Stack) compose(ctx context.Context, io stream.IOStream, args ...string) func() int {
	if ds.DockerExecutable == "" {
		var err error
		ds.DockerExecutable, err = ds.Env.LookPathAbs("docker")
		if err != nil {
			return environment.ExecCommandErrorFunc
		}
	}
	return ds.Env.Exec(ctx, io, ds.Dir, ds.DockerExecutable, append([]string{"compose"}, args...)...)
}

// StackWithResources represents a Stack that can be automatically installed from a set of resources.
// See the [Install] method.
type StackWithResources struct {
	Stack

	// Installable enabled installing several resources from a (potentially embedded) filesystem.
	//
	// The Resources holds these, with appropriate resources specified below.
	// These all refer to paths within the Resource filesystem.
	Resources   fs.FS
	ContextPath string            // the 'docker compose' stack context, containing e.g. 'docker-compose.yml'.
	EnvPath     string            // the '.env' template, will be installed using [unpack.InstallTemplate].
	EnvContext  map[string]string // context when instantiating the '.env' template

	CopyContextFiles []string // Files to copy from the installation context

	MakeDirsPerm fs.FileMode // permission for dirctories, defaults to [environment.DefaultDirCreate]
	MakeDirs     []string    // directories to ensure that exist

	TouchFilesPerm fs.FileMode // permission for new files to touch, defaults to [environment.DefaultFileCreate]
	TouchFiles     []string    // Files to 'touch', i.e. ensure that exist; guaranteed to be run after MakeDirs
}

// InstallationContext is a context to install data in
type InstallationContext map[string]string

// Install installs or updates this stack into the directory specified by stack.Stack().
//
// Installation is non-interactive, but will provide debugging output onto io.
// InstallationContext
func (is StackWithResources) Install(ctx context.Context, progress io.Writer, context InstallationContext) error {
	env := is.Stack.Env
	if is.ContextPath != "" {
		// setup the base files
		if err := unpack.InstallDir(
			env,
			is.Dir,
			is.ContextPath,
			is.Resources,
			func(dst, src string) {
				logging.ProgressF(progress, ctx, "[install] %s\n", dst)
			},
		); err != nil {
			return err
		}
	}

	// configure .env
	envDest := filepath.Join(is.Dir, ".env")
	if is.EnvPath != "" && is.EnvContext != nil {
		logging.ProgressF(progress, ctx, "[config]  %s\n", envDest)
		if err := unpack.InstallTemplate(
			env,
			envDest,
			is.EnvContext,
			is.EnvPath,
			is.Resources,
		); err != nil {
			return err
		}
	}

	// make sure that certain dirs exist
	for _, name := range is.MakeDirs {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		logging.ProgressF(progress, ctx, "[make]    %s\n", dst)
		if is.MakeDirsPerm == fs.FileMode(0) {
			is.MakeDirsPerm = environment.DefaultDirPerm
		}
		if err := env.MkdirAll(dst, is.MakeDirsPerm); err != nil {
			return err
		}
	}

	// copy files from the context!
	for _, name := range is.CopyContextFiles {
		// find the source!
		src, ok := context[name]
		if !ok {
			return errors.Errorf("Missing file from context: %q", src)
		}

		// find the destination!
		dst := filepath.Join(is.Dir, name)

		// copy over file from context
		logging.ProgressF(progress, ctx, "[copy]    %s (from %s)\n", dst, src)
		if err := fsx.CopyFile(ctx, env, dst, src); err != nil {
			return errors.Wrapf(err, "Unable to copy file %s", src)
		}
	}

	// make sure that certain files exist
	for _, name := range is.TouchFiles {
		// find the destination!
		dst := filepath.Join(is.Dir, name)

		logging.ProgressF(progress, ctx, "[touch]   %s\n", dst)
		if err := fsx.Touch(env, dst, is.TouchFilesPerm); err != nil {
			return err
		}
	}

	return nil
}
