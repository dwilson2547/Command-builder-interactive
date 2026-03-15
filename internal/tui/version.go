package tui

// AppVersion is the current application version.
// It is set at build time via:
//
//	go build -ldflags "-X github.com/dwilson2547/command-builder/internal/tui.AppVersion=vX.Y.Z"
//
// The build script (build.sh) handles this automatically.
// Increment the minor version for every change unless otherwise specified.
var AppVersion = "v1.10.0"
