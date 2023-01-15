package config

import (
	"fmt"
	"strings"

	"github.com/tkw1536/goprogram/lib/collection"
)

// This file contains domain related derived configuration values.

// HTTPSEnabled returns if the distillery has HTTPS enabled, and false otherwise.
func (cfg Config) HTTPSEnabled() bool {
	return cfg.CertbotEmail != ""
}

// HostRequirement returns a traefik rule for the given names
func (Config) HostRule(names ...string) string {
	quoted := collection.MapSlice(names, func(name string) string {
		return "`" + name + "`"
	})
	return fmt.Sprintf("Host(%s)", strings.Join(quoted, ","))
}

// HTTPSEnabledEnv returns "true" if https is enabled, and "false" otherwise.
func (cfg Config) HTTPSEnabledEnv() string {
	if cfg.HTTPSEnabled() {
		return "true"
	}
	return "false"
}

// HostFromSlug returns the hostname belonging to a given slug.
// When the slug is empty, returns the default (top-level) domain.
func (cfg Config) HostFromSlug(slug string) string {
	if slug == "" {
		return cfg.DefaultDomain
	}
	return fmt.Sprintf("%s.%s", slug, cfg.DefaultDomain)
}

// DefaultHostRule returns the default traefik hostname rule for this distillery.
// This consists of the [DefaultDomain] as well as [ExtraDomains].
func (cfg Config) DefaultHostRule() string {
	return cfg.HostRule(append([]string{cfg.DefaultDomain}, cfg.SelfExtraDomains...)...)
}

// SlugFromHost returns the slug belonging to the appropriate host.'
//
// When host is a top-level domain, returns "", true.
// When no slug is found, returns "", false.
func (cfg Config) SlugFromHost(host string) (slug string, ok bool) {
	// extract an ':port' that happens to be in the host.
	domain, _, _ := strings.Cut(host, ":")
	domain = TrimSuffixFold(domain, ".wisski") // remove optional ".wisski" ending that is used inside docker

	domainL := strings.ToLower(domain)

	// check all the possible domain endings
	for _, suffix := range append([]string{cfg.DefaultDomain}, cfg.SelfExtraDomains...) {
		suffixL := strings.ToLower(suffix)
		if domainL == suffixL {
			return "", true
		}
		if strings.HasSuffix(domainL, "."+suffixL) {
			return domain[:len(domain)-len(suffix)-1], true
		}
	}

	// no domain found!
	return "", ok
}

func TrimSuffixFold(s string, suffix string) string {
	if len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
