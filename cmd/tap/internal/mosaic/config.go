package mosaic

import (
	"path"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
)

// Config contains paths to the standalone console utils.
type Config struct {
	cmds  string // base directory for commands
	data  string // base directory for data
	types rs.TypeSpecs
}

func (cfg *Config) Scratch(parts ...string) string {
	return filepath.Join(append([]string{cfg.data, "build"}, parts...)...)
}

func (cfg *Config) PathTo(parts ...string) string {
	return filepath.Join(append([]string{cfg.data}, parts...)...)
}

// Rather than creating one big app, for now, tapestry is split into a bunch of separate commands.
func (cfg *Config) Cmd(which string) string {
	return path.Join(cfg.cmds, "bin", which)
}

func (cfg *Config) Stories() string {
	// without the manually added trailing slash, the file open dialog doesnt follow symlink(s) correctly.
	return cfg.PathTo("stories") + "/"
}

// Configure creates a reasonable(?) config based on the developer go path.
func Configure(types rs.TypeSpecs, cmdDir, dataDir string) *Config {
	return &Config{
		types: types,
		cmds:  cmdDir,
		data:  dataDir,
	}
}
