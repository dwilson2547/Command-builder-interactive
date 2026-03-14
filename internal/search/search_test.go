package search

import (
	"testing"

	"github.com/dwilson2547/command-builder/internal/config"
)

func sampleConfigs() []*config.Config {
	return []*config.Config{
		{
			Name:        "default",
			Description: "Built-in",
			Version:     "1.0.0",
			Commands: []config.Command{
				{
					Name:        "openssl",
					Description: "OpenSSL toolkit",
					Options: []config.Option{
						{
							Name:        "print-p12",
							Description: "Print P12 keystore content",
							Template:    "openssl pkcs12 -info -in {{f}}",
						},
						{
							Name:        "generate-rsa-key",
							Description: "Generate RSA private key",
							Template:    "openssl genrsa -out {{f}} 4096",
						},
					},
				},
				{
					Name:        "tar",
					Description: "Tape archive",
					Options: []config.Option{
						{
							Name:        "create-gzip",
							Description: "Create gzip compressed archive",
							Template:    "tar -czvf {{out}} {{src}}",
						},
					},
				},
			},
		},
	}
}

func TestSearchEmpty(t *testing.T) {
	cfgs := sampleConfigs()
	results := Search("", cfgs, Filter{Type: FilterAll})
	if len(results) != 3 {
		t.Errorf("expected 3 results for empty query, got %d", len(results))
	}
}

func TestSearchByCommandName(t *testing.T) {
	cfgs := sampleConfigs()
	results := Search("openssl", cfgs, Filter{Type: FilterAll})
	if len(results) != 2 {
		t.Errorf("expected 2 openssl results, got %d", len(results))
	}
	// All results should be from the openssl command.
	for _, r := range results {
		if r.Command.Name != "openssl" {
			t.Errorf("expected openssl command, got %q", r.Command.Name)
		}
	}
}

func TestSearchByDescription(t *testing.T) {
	cfgs := sampleConfigs()
	results := Search("keystore", cfgs, Filter{Type: FilterAll})
	if len(results) == 0 {
		t.Error("expected at least one result for 'keystore'")
	}
	found := false
	for _, r := range results {
		if r.Option.Name == "print-p12" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected print-p12 option in results for 'keystore'")
	}
}

func TestSearchFilterDefault(t *testing.T) {
	cfgs := append(sampleConfigs(), &config.Config{
		Name: "custom",
		Commands: []config.Command{
			{Name: "mycmd", Options: []config.Option{
				{Name: "myopt", Description: "custom option"},
			}},
		},
	})
	results := Search("", cfgs, Filter{Type: FilterDefault})
	for _, r := range results {
		if r.Config.Name != "default" {
			t.Errorf("FilterDefault returned result from config %q", r.Config.Name)
		}
	}
}

func TestSearchFilterConfig(t *testing.T) {
	cfgs := append(sampleConfigs(), &config.Config{
		Name: "custom",
		Commands: []config.Command{
			{Name: "mycmd", Options: []config.Option{
				{Name: "myopt", Description: "custom option"},
			}},
		},
	})
	results := Search("", cfgs, Filter{Type: FilterConfig, ConfigName: "custom"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Config.Name != "custom" {
		t.Errorf("expected custom config, got %q", results[0].Config.Name)
	}
}

func TestParseQuery(t *testing.T) {
	cases := []struct {
		in         string
		filterType FilterType
		cfgName    string
		terms      string
	}{
		{"/default print", FilterDefault, "", "print"},
		{"/all tar", FilterAll, "", "tar"},
		{"/myconfig ssh", FilterConfig, "myconfig", "ssh"},
		{"openssl rsa", FilterAll, "", "openssl rsa"},
	}
	for _, c := range cases {
		f, terms := ParseQuery(c.in)
		if f.Type != c.filterType {
			t.Errorf("ParseQuery(%q): filter type %v, want %v", c.in, f.Type, c.filterType)
		}
		if f.ConfigName != c.cfgName {
			t.Errorf("ParseQuery(%q): config name %q, want %q", c.in, f.ConfigName, c.cfgName)
		}
		if terms != c.terms {
			t.Errorf("ParseQuery(%q): terms %q, want %q", c.in, terms, c.terms)
		}
	}
}
