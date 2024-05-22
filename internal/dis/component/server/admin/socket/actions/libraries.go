package actions

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/pkglib/stream"
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

func (*InstallColorboxJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_colorbox.sh")
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

func (*InstallDompurifyJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_dompurify.sh")
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

func (*InstallMiradorJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_mirador.sh")
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

func (*InstallIIPMooViewerJS) Act(ctx context.Context, instance *wisski.WissKI, in io.Reader, out io.Writer, params ...string) error {
	return instance.Barrel().Shell(ctx, stream.NewIOStream(out, out, nil), "/runtime/install_iipmooviewer.sh")
}
