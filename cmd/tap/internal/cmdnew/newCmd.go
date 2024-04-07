package cmdnew

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/content"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// called by tap.go
func newCmd(ctx context.Context, _ *base.Command, args []string) (err error) {
	if cnt := len(args); cnt == 0 || cnt > 2 {
		err = fmt.Errorf("%w expected a story for the new story", base.UsageError)
	} else {
		if name := args[0]; name == "sample" {
			if len(args) == 1 {
				err = listSamples()
			} else {
				sample := args[1]
				err = newFromSample(sample, cfg.forceOverwrite)
			}
		} else {
			var title string
			if len(args) == 1 {
				title = "Untitled"
			}
			err = newFromTemplate(name, title, cfg.forceOverwrite)
		}
	}
	return
}

// generate the named story using the default story template
func newFromTemplate(name, title string, force bool) (err error) {
	if n := inflect.Normalize(name); n != strings.ToLower(name) {
		err = fmt.Errorf("%w prefer story names without special characters. For example, maybe %q",
			base.UsageError, n)
	} else {
		out := pathFromName(name)
		if dst, e := createFile(out, force); e != nil {
			err = e
		} else {
			defer dst.Close()
			t := template.Must(template.New("story").Parse(content.DefaultStory))
			if e := t.Execute(dst, content.DefaultDesc{
				Story:  name,
				Title:  title,
				Author: authorName(),
			}); e != nil {
				err = e
			} else {
				fmt.Printf("Created a new story file: %s\n", out)
			}
		}
	}
	return
}

// copy from the named embedded sample to the story directory.
func newFromSample(name string, force bool) (err error) {
	fname := name + files.TellStory.String()
	at := filepath.Join("stories", fname)
	if src, e := content.Sample.Open(at); errors.Is(e, fs.ErrNotExist) {
		err = fmt.Errorf("unknown sample %q. you can list the available samples using `tap new sample`",
			name)
	} else if e != nil {
		err = e
	} else {
		defer src.Close()
		out := pathFromName(name)
		if dst, e := createFile(out, force); e != nil {
			err = e
		} else {
			defer dst.Close()
			if _, e := io.Copy(dst, src); e != nil {
				err = e
			} else {
				fmt.Printf("Wrote sample %q to %s\n", name, out)
			}
		}
	}
	return
}

func listSamples() (err error) {
	if ds, e := fs.ReadDir(content.Sample, "stories"); e != nil {
		err = e
	} else {
		fmt.Printf("The following samples are available: ")
		for i, d := range ds {
			if !d.IsDir() {
				if i > 0 {
					fmt.Print(", ")
				}
				n := d.Name()
				ext := files.TellStory.String()
				fmt.Printf("%q", n[:len(n)-len(ext)])
			}
		}
		fmt.Println(".")
	}
	return
}

func pathFromName(name string) string {
	return filepath.Join(storyDir(), name+files.TellStory.String())
}

func createFile(out string, force bool) (ret *os.File, err error) {
	if _, e := os.Stat(out); e == nil && !force {
		err = fmt.Errorf("The file %q already exists. Remove it first, or use -force to overwrite.", out)
	} else if e != nil && !errors.Is(e, fs.ErrNotExist) {
		err = e
	} else {
		ret, err = os.Create(out)
	}
	return
}

// description of the 'new' command; used by tap.go
var CmdNew = &base.Command{
	Run:       newCmd,
	Flag:      buildFlags(),
	UsageLine: `tap new [-force] name ["title"]`,
	Short:     "create a new story",
	Long:      longDesc(),
}

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	forceOverwrite bool
}{}

// returns a command line parsing object
func buildFlags() (fs flag.FlagSet) {
	fs.BoolVar(&cfg.forceOverwrite, "force", false, "overwrite existing files if needed")
	return
}

func longDesc() string {
	where := storyDir()
	return fmt.Sprintf(
		`Create a new story file, with optional initial title.

The specified name is used as a filename, so single words with 
no special characters are preferred. New story files are created in:
    %s

As a special case, "tap new sample" will list the available sample stories.
And "tap new sample <name of sample>" will create that story.
`, where)
}

func storyDir() (ret string) {
	if home, e := os.UserHomeDir(); e == nil {
		ret = filepath.Join(home, "Documents", "Tapestry", "stories")
	} else if cwd, e := os.Getwd(); e == nil {
		ret = cwd
	} else {
		ret = "."
	}
	return
}
func authorName() (ret string) {
	ret = "Anonymous"
	if u, e := user.Current(); e == nil {
		if n := u.Name; len(n) > 0 {
			ret = n
		} else if n := u.Username; len(n) > 0 {
			ret = n
		}
	}
	return
}
