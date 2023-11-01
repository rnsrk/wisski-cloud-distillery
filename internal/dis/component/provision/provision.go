package provision

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/manager"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/fsx"
)

type Provision struct {
	component.Base
	dependencies struct {
		Instances     *instances.Instances
		Provisionable []component.Provisionable
	}
}

// Flags are flags for a new instance.
type Flags struct {
	// NOTE(twiesing): Any changes here should be reflected in instance_provision.html and remote/api.ts.

	// Slug is the slug of the wisski instance
	Slug string

	// System is information about the system
	System models.System
}

// Profile returns the profile belonging to this provision flags.
func (flags Flags) Profile() manager.Profile {
	// TODO: Actually do something here
	return manager.Profile{}
}

var ErrInstanceAlreadyExists = errors.New("instance with provided slug already exists")

func (pv *Provision) Validate(flags Flags) error {
	// check the slug
	if _, err := pv.dependencies.Instances.IsValidSlug(flags.Slug); err != nil {
		return err
	}
	return nil
}

// Provision provisions a new docker compose instance.
func (pv *Provision) Provision(progress io.Writer, ctx context.Context, flags Flags) (*wisski.WissKI, error) {
	// check that it doesn't already exist
	logging.LogMessage(progress, "Provisioning new WissKI instance %s", flags.Slug)
	if exists, err := pv.dependencies.Instances.Has(ctx, flags.Slug); err != nil || exists {
		return nil, ErrInstanceAlreadyExists
	}

	// log out what we're doing!
	fmt.Fprintf(progress, "%#v", flags)

	// make it in-memory
	instance, err := pv.dependencies.Instances.Create(flags.Slug, flags.System)
	if err != nil {
		return nil, err
	}

	// check that the base directory does not exist
	{
		logging.LogMessage(progress, "Checking that base directory %s does not exist", instance.FilesystemBase)
		exists, err := fsx.Exists(instance.FilesystemBase)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrInstanceAlreadyExists
		}
	}

	// Store in the instances table!
	if err := logging.LogOperation(func() error {
		if err := instance.Bookkeeping().Save(ctx); err != nil {
			return err
		}

		return nil
	}, progress, "Updating bookkeeping database"); err != nil {
		return nil, err
	}

	// create all the resources!
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range pv.dependencies.Provisionable {
			logging.LogMessage(progress, "Provisioning %s resources", pc.Name())
			err := pc.Provision(ctx, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, progress, "Provisioning instance-specific resources"); err != nil {
		return nil, err
	}

	// run the provision script
	if err := logging.LogOperation(func() error {
		return instance.Manager().Provision(ctx, progress, flags.System, flags.Profile())
	}, progress, "Running setup scripts"); err != nil {
		return nil, err
	}

	// start the container!
	logging.LogMessage(progress, "Starting Container")
	if err := instance.Barrel().Stack().Up(ctx, progress); err != nil {
		return nil, err
	}

	// and return the instance
	return instance, nil
}
