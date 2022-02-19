package composer

import (
	"log"
	"os"
	"path"
	"strconv"
)

// Config contains paths to the standalone console utils.
// Rather than creating one big app, for now, tapestry is split into a bunch of separate commands.
type Config struct {
	Assemble string
	Check    string
	Play     string
	Root     string
	Port     int
	Mosaic   string
}

func (cfg *Config) PortString() string {
	return ":" + strconv.Itoa(cfg.Port)
}

func (cfg *Config) Scratch(parts ...string) string {
	return path.Join(append([]string{cfg.Root, "build"}, parts...)...)
}

func (cfg *Config) PathTo(parts ...string) string {
	return path.Join(append([]string{cfg.Root}, parts...)...)
}

// DevConfig creates a reasonable(?) config based on the developer go path.
func DevConfig(base string) *Config {
	bin := "bin"
	var dir string // echo $TMPDIR
	if temp, e := os.MkdirTemp("", "tap"); e != nil {
		log.Fatal(e)
	} else {
		dir = temp
	}
	a, cfg, p := "asm", "check", "play"
	return &Config{
		Assemble: path.Join(base, bin, a),
		Check:    path.Join(base, bin, cfg),
		Play:     path.Join(base, bin, p),
		Root:     dir,
		Port:     3000,
	}
}
