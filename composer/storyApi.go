package composer

import (
	"context"
	"io"
	"log"
	"os/exec"
	"path"
	"strings"

	"git.sr.ht/~ionous/tapestry/web"
	"github.com/ionous/errutil"
)

// uses the command line tool versions for now....
func tempTest(ctx context.Context, file string, in io.Reader) (err error) {
	cfg := ctx.Value(configKey).(*web.Config)
	base := cfg.PathTo("stories")
	if !strings.HasPrefix(file, base) {
		err = errutil.New("unexpected path", file, "from", base)
	} else {
		// note: .Split keeps a trailing slash, .Dir does not.
		dir, _ := path.Split(file[len(base)+1:])
		const shared = "shared/"
		const stories = "stories/"
		// we'll always include the shared files in our build
		src := cfg.PathTo(stories, shared)
		if strings.HasPrefix(dir, shared) {
			dir = shared
		} else {
			// get the first part of the name -- that's the project name
			i := strings.Index(dir, "/")
			dir = dir[0:i] // the project relative dir
			src += "," + cfg.PathTo(stories, dir)
		}
		// src is now one or two absolute paths to project directories
		// dir is a relative dir
		if playFile, e := runAsm(ctx, cfg, src, dir); e != nil {
			log.Println(tab, "Assembly error", cfg.Cmd("asm"), exitError(e))
			err = e
		} else if e := runCheck(ctx, cfg, playFile); e != nil {
			log.Println(tab, "Check error", cfg.Cmd("check"), exitError(e))
			err = e
		}
	}
	return
}

const tab = '\t'

func runAsm(ctx context.Context, cfg *web.Config, ephFile, path string) (ret string, err error) {
	log.Println("Assembling", ephFile+"...")
	inFile, playFile := ephFile, cfg.Scratch(path, "play.db")
	log.Println(">", cfg.Cmd("asm"), "-in", inFile, "-out", playFile)
	assembled, e := exec.CommandContext(ctx, cfg.Cmd("asm"), "-in", inFile, "-out", playFile).CombinedOutput()
	if e != nil {
		err = e
	} else {
		ret = playFile
	}
	logBytes(assembled)
	return

}
func runCheck(ctx context.Context, cfg *web.Config, playFile string) (err error) {
	log.Println("Checking", playFile+"...")
	log.Println(">", cfg.Cmd("check"), "-in", playFile)
	checked, e := exec.CommandContext(ctx, cfg.Cmd("check"), "-in", playFile).CombinedOutput()
	if e != nil {
		err = e
	}
	logBytes(checked)
	return
}

func logBytes(b []byte) {
	if s := strings.Trim(string(b), "\n"); len(s) > 0 {
		log.Println(s)
	}
}

func exitError(e error) (ret string) {
	if x, ok := e.(*exec.ExitError); ok {
		ret = string(x.Stderr)
	} else {
		ret = "unknown."
	}
	return
}
