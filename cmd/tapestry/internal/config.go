package tap

import (
	"path"
	"path/filepath"
)

// Config contains paths to the standalone console utils.
type Config struct {
	cmds string // base directory for commands
	data string // base directory for data
	prod string // if this exists: a packaged set of frontend assets
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

// empty string if not production
func (cfg *Config) Prod() string {
	return cfg.prod
}

// DevConfig creates a reasonable(?) config based on the developer go path.
func DevConfig(cmdDir, dataDir string) *Config {
	return &Config{
		cmds: cmdDir,
		data: dataDir,
	}
}
