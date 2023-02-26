package config

import (
	"crypto/rand"

	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/tkw1536/pkglib/password"
)

// NewPassword returns a new password using the password settings from this configuration
func (cfg Config) NewPassword() (string, error) {
	return password.Generate(rand.Reader, cfg.PasswordLength, passwordx.Charset)
}
