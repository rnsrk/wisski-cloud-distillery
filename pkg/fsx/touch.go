package fsx

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Touch touches a file.
// It is similar to the unix 'touch' command.
//
// If the file does not exist exists, it is created using [env.Create].
// If the file does exist, it's access and modification times are updated to the current time.
func Touch(env environment.Environment, path string) error {
	_, err := env.Stat(path)
	switch {
	case environment.IsNotExist(err):
		f, err := env.Create(path, environment.DefaultFilePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	case err != nil:
		return err
	default:
		now := time.Now().Local()
		return env.Chtimes(path, now, now)
	}
}
