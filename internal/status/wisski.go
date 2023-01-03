package status

import (
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
)

// WissKI provides information about a single WissKI
type WissKI struct {
	Time time.Time // Time this info was built

	Slug string // slug
	URL  string // complete URL, including http(s)

	Locked bool // Is this instance currently locked?

	// Information about the running instance
	Running     bool
	LastRebuild time.Time
	LastUpdate  time.Time
	LastCron    time.Time

	// Statistics of the WissKI
	Statistics Statistics

	// List of backups made
	Snapshots []models.Export

	// List of SSH Keys
	SSHKeys []string

	// WissKI content information
	NoPrefixes   bool              // TODO: Move this into the database
	Prefixes     []string          // list of prefixes
	Pathbuilders map[string]string // all the pathbuilders
	Users        []User            // all the known users
	Grants       []models.Grant
}

// Statistics holds statistics generated by the WissKI module
type Statistics struct {
	Activity struct {
		MostVisited string `json:"mostVisited"`
		PageVisits  []struct {
			URL    string `json:"url"`
			Visits int    `json:"visits"`
		} `json:"pageVisits"`
		TotalEditsLastWeek int `json:"totalEditsLastWeek"`
	} `json:"activity"`
	Bundles     BundleStatistics `json:"bundles"`
	Triplestore struct {
		Graphs []struct {
			URI   string `json:"uri"`
			Count int    `json:"triples"`
		} `json:"graphStatistics"`
		Total int `json:"totalTriples"`
	} `json:"triplestore"`
	Users struct {
		LastLogin  string `json:"lastLogin"`
		TotalUsers int    `json:"totalUsers"`
	} `json:"users"`
}

type BundleStatistics struct {
	Bundles []struct {
		Label       string `json:"label"`
		MachineName string `json:"machineName"`

		Count int `json:"entities"`

		LastEdit phpx.Timestamp `json:"lastEdit"`

		MainBundle phpx.Boolean `json:"mainBundle"`
	} `json:"bundleStatistics"`
	TotalBundles     int `json:"totalBundles"`
	TotalMainBundles int `json:"totalMainBundles"`
}

type LastEdit struct {
	Time  time.Time
	Valid bool
}

// LastEdit returns the last time any bundle was edited, and if any edit was bigger than the reference time
func (bs BundleStatistics) LastEdit() (le LastEdit) {
	for _, bundle := range bs.Bundles {
		time := bundle.LastEdit.Time()
		// skip invalid times
		if time.Unix() <= 0 {
			continue
		}
		if time.After(le.Time) {
			le.Valid = true
			le.Time = time
		}
	}
	return
}

func (bs BundleStatistics) Summary() string {
	var totalCount int
	for _, bundle := range bs.Bundles {
		totalCount += bundle.Count
	}
	if totalCount == 0 {
		return ""
	}

	entitySubject := "Entities"
	if totalCount == 1 {
		entitySubject = "Entity"
	}

	bundleSubject := "Bundles"
	if len(bs.Bundles) == 1 {
		bundleSubject = "Bundle"
	}

	return fmt.Sprintf("%d %s in %d %s", totalCount, entitySubject, len(bs.Bundles), bundleSubject)
}
