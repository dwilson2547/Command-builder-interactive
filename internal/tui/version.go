package tui

// AppVersion is the current application version.
// It is set at build time via:
//
//	go build -ldflags "-X github.com/dwilson2547/command-builder/internal/tui.AppVersion=vX.Y.Z"
//
// The build script (build.sh) handles this automatically.
// Increment the minor version for every change unless otherwise specified.
var AppVersion = "v1.39.0"

// AppDisplayName is the application name shown in the header.
// It defaults to "Command Builder" and is updated from user settings via ApplyTheme.
var AppDisplayName = "Command Builder"
