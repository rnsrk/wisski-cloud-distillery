// Package core implements the core of the WissKI Distillery and the wdcli executable.
// It does not depend on any other packages.
package core

import _ "embed"

// BaseDirectoryDefault is the default deploy directory to load the distillery from.
const BaseDirectoryDefault = "/var/www/deploy"

// Executable is the name of the 'wdcli' executable.
// It should be located inside the deployment directory.
const Executable = "wdcli"

// ConfigFile is the name of the config file.
// It should be located inside the deployment directory.
const ConfigFile = ".env"

// OverridesJSON is the name of the json overrides file.
// It should be located inside the deployment directory.
const OverridesJSON = "overrides.json"

// DefaultOverridesJSON contains a template for a new 'overrides.json' file
//go:embed bootstrap/overrides.json
var DefaultOverridesJSON []byte

// AuthorizedKeys contains the default name for the 'global_authorized_keys' file
const AuthorizedKeys = "authorized_keys"

// DefaultAuthorizedKeys contains a template for a new 'global_authorized_keys' file
//go:embed bootstrap/global_authorized_keys
var DefaultAuthorizedKeys []byte

// PrefixConfig is the name for the global resolver prefix configuration.
// It should be found within the prefix component directory.
const PrefixConfig = "prefix.cfg"
