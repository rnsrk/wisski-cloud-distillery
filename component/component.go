// Package component holds the main abstraction for components.
package component

import (
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// Component represents a logical subsystem of the distillery.
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
// Furthermore, a component is typically instantiated using a call on the ["distillery.Distillery"] struct.
//
// Each Component should make use of [ComponentBase] for sane defaults.
//
// For example, the web.Web component lives in the web package and can be created like:
//
//	var dis Distillery
//  web := dis.Web()
type Component interface {
	// Name returns the name of this component.
	// It should correspond to the appropriate subpackage.
	Name() string

	// Path returns the path this component is installed at.
	// By convention it is /var/www/deploy/core/${Name()}
	Path() string

	// Stack can be used to gain access to the "docker compose" stack.
	//
	// This should internally call
	Stack() stack.Installable

	// Context returns a new InstallationContext to be used during installation from the command line.
	// Typically this should just pass through the parent, but might perform other tasks.
	Context(parent stack.InstallationContext) stack.InstallationContext
}

// ComponentBase implements base functionality for a component
type ComponentBase struct {
	Dir string // Dir is the directory this component lives in

	Config *config.Config // Config is the configuration of the underlying distillery
}

// Path returns the path to this component
func (cb ComponentBase) Path() string {
	return cb.Dir
}

// Context passes through the parent context
func (ComponentBase) Context(parent stack.InstallationContext) stack.InstallationContext {
	return parent
}

// MakeStack registers the Installable as a stack
func (cb ComponentBase) MakeStack(stack stack.Installable) stack.Installable {
	stack.Dir = cb.Dir
	return stack
}