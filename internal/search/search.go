// Package search provides fuzzy search over loaded command configs.
package search

import (
	"sort"
	"strings"

	"github.com/dwilson2547/command-builder/internal/config"
)

// FilterType controls which configs are searched.
type FilterType int

const (
	FilterAll     FilterType = iota // search everything
	FilterDefault                   // only the "default" built-in config
	FilterConfig                    // a specific named config
)

// Filter holds the active search filter.
type Filter struct {
	Type       FilterType
	ConfigName string // only used when Type == FilterConfig
}

// SearchResult represents a single matching command option.
type SearchResult struct {
	Config  *config.Config
	Command *config.Command
	Option  *config.Option
	Score   int
}

// ParseQuery extracts a Filter and the raw search terms from a query string.
// Rules:
//
//	"/default <terms>" → FilterDefault
//	"/all <terms>"     → FilterAll
//	"/<name> <terms>"  → FilterConfig{name}
//	"<terms>"          → FilterAll
func ParseQuery(query string) (Filter, string) {
	q := strings.TrimSpace(query)
	if strings.HasPrefix(q, "/") {
		parts := strings.SplitN(q[1:], " ", 2)
		modifier := strings.ToLower(parts[0])
		rest := ""
		if len(parts) == 2 {
			rest = parts[1]
		}
		switch modifier {
		case "all":
			return Filter{Type: FilterAll}, rest
		case "default":
			return Filter{Type: FilterDefault}, rest
		default:
			return Filter{Type: FilterConfig, ConfigName: modifier}, rest
		}
	}
	return Filter{Type: FilterAll}, q
}

// Search returns results matching query across the given configs, applying
// the filter. Results are sorted descending by score.
// The caller is responsible for extracting the filter via ParseQuery; the
// query string is used here only to derive search terms.
func Search(query string, configs []*config.Config, filter Filter) []SearchResult {
	// Extract terms only; trust the caller's filter.
	_, terms := ParseQuery(query)

	var results []SearchResult

	for _, cfg := range configs {
		// Apply filter.
		switch filter.Type {
		case FilterDefault:
			if cfg.Name != "default" {
				continue
			}
		case FilterConfig:
			if cfg.Name != filter.ConfigName {
				continue
			}
		}

		for ci := range cfg.Commands {
			cmd := &cfg.Commands[ci]
			for oi := range cmd.Options {
				opt := &cmd.Options[oi]
				score := scoreMatch(terms, cfg, cmd, opt)
				if score > 0 {
					results = append(results, SearchResult{
						Config:  cfg,
						Command: cmd,
						Option:  opt,
						Score:   score,
					})
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// scoreMatch returns a relevance score for the given option against terms.
// Returns 0 if the option should be excluded from results.
func scoreMatch(query string, cfg *config.Config, cmd *config.Command, opt *config.Option) int {
	if query == "" {
		return 1 // all results when query is empty
	}

	terms := strings.Fields(strings.ToLower(query))
	score := 0

	// Build a searchable corpus for this option.
	cmdName := strings.ToLower(cmd.Name)
	optName := strings.ToLower(opt.Name)
	optDesc := strings.ToLower(opt.Description)
	cmdDesc := strings.ToLower(cmd.Description)
	cfgName := strings.ToLower(cfg.Name)

	// Build a lower-cased list of tags for this option.
	var optTags []string
	for _, t := range opt.Tags {
		optTags = append(optTags, strings.ToLower(t))
	}

	for _, term := range terms {
		termScore := 0

		// Exact match on command name scores highest.
		if cmdName == term {
			termScore += 100
		} else if strings.HasPrefix(cmdName, term) {
			termScore += 60
		} else if strings.Contains(cmdName, term) {
			termScore += 30
		}

		// Option name.
		if optName == term {
			termScore += 80
		} else if strings.HasPrefix(optName, term) {
			termScore += 50
		} else if strings.Contains(optName, term) {
			termScore += 25
		}

		// Description matches.
		if strings.Contains(optDesc, term) {
			termScore += 20
		}
		if strings.Contains(cmdDesc, term) {
			termScore += 10
		}
		if strings.Contains(cfgName, term) {
			termScore += 5
		}

		// Tag matches – exact tag match scores like an option name match.
		for _, tag := range optTags {
			if tag == term {
				termScore += 80
				break
			} else if strings.HasPrefix(tag, term) {
				termScore += 50
				break
			} else if strings.Contains(tag, term) {
				termScore += 25
				break
			}
		}

		// Fuzzy: each rune of the term must appear in order somewhere.
		if termScore == 0 {
			// Include tags in fuzzy corpus.
			tagCorpus := strings.Join(optTags, " ")
			fuzzyCorpus := optDesc + " " + cmdName + " " + optName
			if tagCorpus != "" {
				fuzzyCorpus += " " + tagCorpus
			}
			if fuzzyContains(fuzzyCorpus, term) {
				termScore += 5
			}
		}

		if termScore == 0 {
			// This term matched nothing – exclude the result.
			return 0
		}
		score += termScore
	}

	return score
}

// fuzzyContains returns true if all runes of sub appear in order in s.
func fuzzyContains(s, sub string) bool {
	si := 0
	for _, r := range sub {
		found := false
		for si < len(s) {
			if rune(s[si]) == r {
				si++
				found = true
				break
			}
			si++
		}
		if !found {
			return false
		}
	}
	return true
}
