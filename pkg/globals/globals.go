package globals

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Globals struct {
	LogLevel string      `short:"l" help:"Set the logging level (${log_levels})" enum:"${log_levels}" default:"info"`
	Config   string      `short:"c" help:"The path to the config file" type:"path" default:"${config_file}"`
	Version  VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"]) // nolint
	app.Exit(0)
	return nil
}
