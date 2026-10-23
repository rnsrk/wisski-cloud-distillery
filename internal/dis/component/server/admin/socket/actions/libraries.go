//spellchecker:words actions
package actions

//spellchecker:words context github wisski distillery internal component auth scopes pkglib stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/pkglib/stream"
)

// installing additional javascript libraries

type InstallColorboxJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallColorboxJS)(nil)
)

func (*InstallColorboxJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-colorbox-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallColorboxJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().BashScript(ctx, stream.NewIOStream(out, out, nil), "/bin/bash", "/runtime/install_colorbox.sh"); err != nil {
		return nil, fmt.Errorf("failed to install colorbox: %w", err)
	}
	return nil, nil
}

type InstallDompurifyJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallDompurifyJS)(nil)
)

func (*InstallDompurifyJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-dompurify-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallDompurifyJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().BashScript(ctx, stream.NewIOStream(out, out, nil), "/bin/bash", "/runtime/install_dompurify.sh"); err != nil {
		return nil, fmt.Errorf("failed to install dompurify: %w", err)
	}
	return nil, nil
}

type InstallMiradorJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallMiradorJS)(nil)
)

func (*InstallMiradorJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-mirador-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallMiradorJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().BashScript(ctx, stream.NewIOStream(out, out, nil), "/bin/bash", "/runtime/install_mirador.sh"); err != nil {
		return nil, fmt.Errorf("failed to install mirador: %w", err)
	}
	return nil, nil
}

type InstallIIPMooViewerJS struct {
	component.Base
}

var (
	_ WebsocketInstanceAction = (*InstallIIPMooViewerJS)(nil)
)

func (*InstallIIPMooViewerJS) Action() InstanceAction {
	return InstanceAction{
		Action: Action{
			Name:      "install-iipmooviewer-js",
			Scope:     scopes.ScopeUserAdmin,
			NumParams: 0,
		},
	}
}

func (*InstallIIPMooViewerJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) (any, error) {
	if err := instance.Barrel().BashScript(ctx, stream.NewIOStream(out, out, nil), "/bin/bash", "/runtime/install_iipmooviewer.sh"); err != nil {
		return nil, fmt.Errorf("failed to install iipmooviewer: %w", err)
	}
	return nil, nil
}
