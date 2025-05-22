//spellchecker:words manager
package manager

//spellchecker:words maps slices github wisski distillery internal ingredient barrel composer drush system bookkeeping extras
import (
	"maps"
	"slices"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/system"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
)

// Manager manages a profile applied to specific WissKI instances.
type Manager struct {
	ingredient.Base
	dependencies struct {
		Barrel      *barrel.Barrel
		Bookkeeping *bookkeeping.Bookkeeping

		SystemManager *system.SystemManager

		Composer *composer.Composer
		Drush    *drush.Drush

		Adapters *extras.Adapters
		Settings *extras.Settings
	}
}

// profiles contains the list of default profiles.
var (
	defaultProfile = "Drupal 11"
	profiles       = map[string]Profile{
		"Drupal 9": {
			Description: "Legacy Version of Drupal",

			Drupal: "^9",
			WissKI: "",
			InstallModules: []string{
				"drupal/inline_entity_form:^3.0@RC",
				"drupal/imagemagick",
				"drupal/image_effects",
				"drupal/colorbox",
			},
			EnableModules: []string{
				"drupal/devel:^4.1",
				"drupal/geofield:^1.40",
				"drupal/geofield_map:^2.85",
				"drupal/imce:^2.4",
				"drupal/remove_generator:^2.0",
			},
		},
		"Drupal 10": {
			Description: "Legacy Version Of Drupal",
			Description: "Drupal 10 with default packages",

			Drupal: "^10",
			WissKI: "4.x-dev@dev",
			InstallModules: []string{
				"drupal/colorbox",
				"drupal/conditional_fields:4.x-dev@dev",
				"drupal/devel:^5.0",
				"drupal/ds:^3.22",
				"drupal/field_group:3.x-dev@dev",
				"drupal/geofield:^1.56",
				"drupal/geofield_map:^3.0",
				"kint-php/kint:^5",
				"drupal/leaflet:^10.2",
				"drupal/imagemagick",
				"drupal/image_effects",
				"drupal/imce:^3.0",
				"drupal/inline_entity_form:^3.0@RC",
			},
			EnableModules: []string{
				"devel",
				"geofield",
				"geofield_map",
				"imce",
				"ds",
				"leaflet",
				"inline_entity_form",
				"colorbox",
				"imagemagick",
				"image_effects",
				"imce",
			},
		},
		"Drupal 11": {
			Description: "Drupal 11 with colorbox, devel, ds, geofield, geofield_map, kint, leaflet, imagemagick, image_effects, imce and inline_entity_form",

			Drupal: "^11",
			WissKI: "4.x-dev@dev",
			InstallModules: []string{
				"drupal/colorbox",
				"drupal/conditional_fields:4.x-dev@dev",
				"drupal/devel:^5",
				"drupal/ds:^3",
				"drupal/field_group:^4.0",
				"drupal/geofield_map:^11",
				"drupal/geofield:^1",
				"drupal/image_effects:^4",
				"drupal/imagemagick:^4",
				"drupal/imce:^3",
				"drupal/inline_entity_form:^3.0@RC",
				"drupal/leaflet:^10",
				"kint-php/kint:^5",
			},
			EnableModules: []string{
				"colorbox",
				"conditional_fields",
				"devel",
				"ds",
				"field_group",
				"geofield_map",
				"geofield",
				"image_effects",
				"imagemagick",
				"imce",
				"inline_entity_form",
				"leaflet",
			},
		},
	}
)

// TODO: All of these should move to the config

func LoadDefaultProfile() Profile {
	return LoadProfile(DefaultProfile())
}

func Profiles() map[string]Profile {
	return maps.Clone(profiles)
}

func LoadProfile(name string) Profile {
	return profiles[name]
}

func HasProfile(name string) bool {
	_, ok := profiles[name]
	return ok
}

func DefaultProfile() string {
	return defaultProfile
}

// Profile represents a profile applied to a WissKI instance of the Distillery.
type Profile struct {
	// Description is a human-readable description for this profile.
	// It is only used by the frontend.
	Description string

	Drupal string // Version of Drupal to use
	WissKI string // Version of WissKI to use

	InstallModules []string // Modules to be installed (but not neccessarily enabled)
	EnableModules  []string // Modules to be installed and enabled
}

// Apply copies over defaults from the other profile to this one.
// If a field is already set, no defaults are copied.
func (profile *Profile) Apply(other Profile) {
	if profile.Drupal == "" {
		profile.Drupal = other.Drupal
	}
	if profile.WissKI == "" {
		profile.WissKI = other.WissKI
	}
	if profile.InstallModules == nil {
		profile.InstallModules = slices.Clone(other.InstallModules)
	}
	if profile.EnableModules == nil {
		profile.EnableModules = slices.Clone(profile.EnableModules)
	}
}

// ApplyDefaults loads some set of defaults.
// If all fields are set, no defaults are applied.
func (profile *Profile) ApplyDefaults() {
	profile.Apply(profiles[defaultProfile])
}
