package env

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// TODO: Move everything into specific subpackages

// Stacks returns the Stacks of this distillery
func (dis *Distillery) Stacks() []stack.Installable {
	// TODO: Do we want to cache these stacks?
	return []stack.Installable{
		dis.WebStack(),
		dis.SelfStack(),
		dis.ResolverStack(),
		dis.SSHStack(),
		dis.TriplestoreStack(),
		dis.SQLStack(),
	}
}

// asCoreStack treats the provided stack as a core component of this distillery.
func (dis *Distillery) asCoreStack(name string, stack stack.Installable) stack.Installable {
	stack.Dir = filepath.Join(dis.Config.DeployRoot, "core", name)

	stack.ContextResource = filepath.Join("resources", "compose", name)
	stack.EnvFileResource = filepath.Join("resources", "templates", "docker-env", name)

	return stack
}
