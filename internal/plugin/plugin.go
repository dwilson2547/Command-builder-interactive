// Package plugin handles importing command configs from remote URLs.
package plugin

import (
	"github.com/dwilson2547/command-builder/internal/config"
)

// ImportFromURL fetches a YAML config from a URL and adds it to the manager.
// This is a convenience wrapper around config.Manager.ImportConfigFromURL.
func ImportFromURL(mgr *config.Manager, url string) (*config.Config, error) {
	return mgr.ImportConfigFromURL(url)
}
