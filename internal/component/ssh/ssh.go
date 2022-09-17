package ssh

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
)

type SSH struct {
	component.ComponentBase
}

func (SSH) Name() string {
	return "ssh"
}

//go:embed all:stack
var resources embed.FS

func (ssh SSH) Stack() component.StackWithResources {
	return ssh.ComponentBase.MakeStack(component.StackWithResources{
		Resources:   resources,
		ContextPath: "stack",
	})
}
