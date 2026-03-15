package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromBytes(t *testing.T) {
	yaml := `
name: "test"
description: "Test config"
version: "1.0.0"
commands:
  - name: "echo"
    description: "Echo command"
    options:
      - name: "simple"
        description: "Simple echo"
        template: "echo {{message}}"
        inputs:
          - name: "message"
            type: "string"
            description: "Message to echo"
            required: true
`
	cfg, err := LoadConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("LoadConfigFromBytes: %v", err)
	}
	if cfg.Name != "test" {
		t.Errorf("expected name 'test', got %q", cfg.Name)
	}
	if len(cfg.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(cfg.Commands))
	}
	if len(cfg.Commands[0].Options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(cfg.Commands[0].Options))
	}
	if len(cfg.Commands[0].Options[0].Inputs) != 1 {
		t.Fatalf("expected 1 input, got %d", len(cfg.Commands[0].Options[0].Inputs))
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	cfg := &Config{
		Name:        "save-test",
		Description: "Save test",
		Version:     "1.0.0",
		Commands: []Command{
			{
				Name:        "testcmd",
				Description: "A test command",
				Options: []Option{
					{
						Name:        "opt1",
						Description: "Option 1",
						Template:    "testcmd --foo {{foo}}",
						Inputs: []Input{
							{
								Name:        "foo",
								Type:        "string",
								Description: "Foo value",
								Required:    true,
							},
						},
					},
				},
			},
		},
	}

	tmp, err := os.CreateTemp("", "cmd-builder-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := SaveConfig(cfg, tmp.Name()); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	loaded, err := LoadConfig(tmp.Name())
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if loaded.Name != cfg.Name {
		t.Errorf("expected name %q, got %q", cfg.Name, loaded.Name)
	}
	if len(loaded.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(loaded.Commands))
	}
}

func TestManagerAddAndDelete(t *testing.T) {
	// Use empty default data so no embedded config is loaded.
	mgr, err := NewManager(nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	cfg := &Config{
		Name:        "mgr-test-" + t.Name(),
		Description: "Manager test",
		Version:     "1.0.0",
	}
	if err := mgr.AddConfig(cfg); err != nil {
		t.Fatalf("AddConfig: %v", err)
	}
	if mgr.GetConfig(cfg.Name) == nil {
		t.Error("expected to find added config")
	}
	if err := mgr.DeleteConfig(cfg.Name); err != nil {
		t.Fatalf("DeleteConfig: %v", err)
	}
	if mgr.GetConfig(cfg.Name) != nil {
		t.Error("expected config to be deleted")
	}
}

func TestSanitizeName(t *testing.T) {
	cases := []struct{ in, out string }{
		{"hello", "hello"},
		{"hello world", "hello-world"},
		{"foo/bar", "foo-bar"},
		{"a:b*c?d", "a-b-c-d"},
	}
	for _, c := range cases {
		got := sanitizeName(c.in)
		if got != c.out {
			t.Errorf("sanitizeName(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}
