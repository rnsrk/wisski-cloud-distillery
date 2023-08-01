package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
)

func (sockets *Sockets) Actions() ActionMap {
	return map[string]Action{
		// generic actions
		"backup": sockets.Generic(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, sockets *Sockets, in io.Reader, out io.Writer, params ...string) error {
			return sockets.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: nil,

					StagingOnly: false,
				},
			)
		}),
		"provision": sockets.Generic(scopes.ScopeUserAdmin, "", 1, func(ctx context.Context, sockets *Sockets, in io.Reader, out io.Writer, params ...string) error {
			// read the flags of the instance to be provisioned
			var flags provision.Flags
			if err := json.Unmarshal([]byte(params[0]), &flags); err != nil {
				return err
			}

			instance, err := sockets.Dependencies.Provision.Provision(
				out,
				ctx,
				flags,
			)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "URL:      %s\n", instance.URL().String())
			fmt.Fprintf(out, "Username: %s\n", instance.DrupalUsername)
			fmt.Fprintf(out, "Password: %s\n", instance.DrupalPassword)

			return nil
		}),

		// instance-specific actions!

		"snapshot": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, socket *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return socket.Dependencies.Exporter.MakeExport(
				ctx,
				out,
				exporter.ExportTask{
					Dest:     "",
					Instance: instance,

					StagingOnly: false,
				},
			)
		}),
		"rebuild": sockets.Instance(scopes.ScopeUserAdmin, "", 1, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			// read the flags of the instance to be provisioned
			var system models.System
			if err := json.Unmarshal([]byte(params[0]), &system); err != nil {
				return err
			}
			return instance.SystemManager().Apply(ctx, out, system, true)
		}),
		"update": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Composer().Update(ctx, out)
		}),
		"cron": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, str io.Writer, params ...string) error {
			return instance.Drush().Cron(ctx, str)
		}),
		"start": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Up(ctx, out)
		}),
		"stop": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, _ *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return instance.Barrel().Stack().Down(ctx, out)
		}),
		"purge": sockets.Instance(scopes.ScopeUserAdmin, "", 0, func(ctx context.Context, sockets *Sockets, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
			return sockets.Dependencies.Purger.Purge(ctx, out, instance.Slug)
		}),
		"never": sockets.Generic(scopes.ScopeNever, "", 0, func(ctx context.Context, sockets *Sockets, in io.Reader, out io.Writer, params ...string) error {
			panic("never called")
		}),
	}
}
