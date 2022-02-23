package web

import (
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

// Config contains paths to the standalone console utils.
type Config struct {
	cmds string // base directory for commands
	data string // base directory for data
}

func (cfg *Config) Scratch(parts ...string) string {
	return path.Join(append([]string{cfg.data, "build"}, parts...)...)
}

func (cfg *Config) PathTo(parts ...string) string {
	return path.Join(append([]string{cfg.data}, parts...)...)
}

// Rather than creating one big app, for now, tapestry is split into a bunch of separate commands.
func (cfg *Config) Cmd(which string) string {
	return path.Join(cfg.cmds, "bin", which)
}

// DevConfig creates a reasonable(?) config based on the developer go path.
func DevConfig(cmdDir, dataDir string) *Config {
	if len(dataDir) == 0 {
		if temp, e := os.MkdirTemp("", "tap"); e != nil {
			log.Fatal(e)
		} else {
			dataDir = temp
		}
	}
	return &Config{
		cmds: cmdDir,
		data: dataDir,
	}
}

func Endpoint(port int, parts ...string) (ret string) {
	ret = ":" + strconv.Itoa(port)
	if len(parts) > 0 {
		u := url.URL{Scheme: "http", Host: parts[0] + ret, Path: strings.Join(parts[1:], "/")}
		ret = u.String() + "/"
	}
	return
}
