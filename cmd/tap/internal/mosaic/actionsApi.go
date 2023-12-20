package mosaic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/useraction"
	"github.com/ionous/errutil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// to open dialogs wails needs its own context:
// it doesn't really follow the standard golang context usag,:
// its a map of storage for various wails values.
// ( the per request context is derived from it. )
// by copying it -- we might miss some later things added to it;
// but mostly its going to be okay.
// note: its not valid until startup has happened.
type Workspace struct {
	Context    context.Context
	hasStarted bool
}

// callback from wails
func (ws *Workspace) Startup(ctx context.Context) {
	ws.Context = ctx
	// alt: could check the .Value("frontend")
	// or wait to attach the api until wails is started?
	// ( not sure if the mux interface allows that )
	ws.hasStarted = true
}

// return the path starting from base
func chopPath(base, path string) (ret string, err error) {
	// the dialog code returns an expanded path, so we need to do the same
	if sym, e := filepath.EvalSymlinks(base); e != nil {
		err = e
	} else if !strings.HasPrefix(path, sym) {
		err = errutil.New("file should be beneath the directory:", base)
	} else {
		ret = path[len(sym)+1:]
	}
	return
}

func ActionsApi(opt *Config, ws *Workspace) web.Resource {
	base := opt.Stories()
	// fix: who could ever call "close()"
	dispatch := useraction.Start(func(action string) (ret string, err error) {
		if !ws.hasStarted {
			err = errutil.New("app statup pending")
		} else {
			switch action {
			case "new":
				if path, e := runtime.SaveFileDialog(ws.Context, runtime.SaveDialogOptions{
					DefaultDirectory: base,
					Title:            "Create new story",
					Filters: []runtime.FileFilter{{
						DisplayName: "Story files (*.if,*.tell)",
						Pattern:     "*.if;*.tell", // semicolon separated list of extensions, EG: "*.jpg;*.png"
					}},
					CanCreateDirectories: true,
					// DefaultFilename            string
					// ShowHiddenFiles            bool
					// TreatPackagesAsDirectories bool
				}); e != nil {
					err = e
				} else if rel, e := chopPath(base, path); e != nil {
					err = e
				} else {
					blank := story.StoryFile{StoryStatements: []story.StoryStatement{}}
					if content, e := story.Encode(&blank); e != nil {
						err = e
					} else if e := files.FormattedSave(filepath.Join(base, rel), content, true); e != nil {
						err = e
					} else {
						ret = rel
					}
				}
			case "open":
				if p, e := runtime.OpenFileDialog(ws.Context, runtime.OpenDialogOptions{
					DefaultDirectory: base,
					Title:            "Create new story",
					Filters: []runtime.FileFilter{{
						DisplayName: "Story files (*.if,*.tell)",
						Pattern:     "*.if;*.tell", // semicolon separated list of extensions, EG: "*.jpg;*.png"
					}},
					CanCreateDirectories: true,
					// DefaultFilename            string
					// ShowHiddenFiles            bool
					// TreatPackagesAsDirectories bool
				}); e != nil {
					err = e
				} else if rel, e := chopPath(base, p); e != nil {
					err = e
				} else {
					ret = rel
				}
			default:
				err = errutil.New("unknown action", action)
			}
		}
		return
	})

	return &web.Wrapper{
		Finds: func(str string) (ret web.Resource) {
			// if there are no more slashes after this then:
			// we are either trying to start a new action ( via post )
			// or poll for an existing one ( via get )
			if parts := strings.Split(str, "/"); len(parts) == 1 {
				name := parts[0]
				// if the name doesnt look like an action token then post it as an action.
				if t := useraction.ReadToken(name); !t.Valid() {
					ret = &web.Wrapper{
						Posts: func(_ context.Context, r io.Reader, w http.ResponseWriter) (err error) {
							if res, e := dispatch.Post(name); e == nil {
								// if the token is valid, the value isnt ready.
								if res.Token.Valid() {
									writeMap(w, map[string]string{
										"token": res.Token.String(),
									})
								} else {
									writeMap(w, map[string]string{
										"value": res.Value,
									})
								}
							} else if status, ok := e.(web.Status); ok {
								http.Error(w, e.Error(), int(status))
							} else {
								err = e
							}
							return
						},
					}
				} else {
					// the name looked like an action token, so poll using get:
					ret = &web.Wrapper{
						Gets: func(_ context.Context, w http.ResponseWriter) (err error) {
							if res, e := dispatch.Get(t); e == nil {
								writeMap(w, map[string]string{
									"value": res,
								})
							} else if status, ok := e.(web.Status); !ok {
								err = e
							} else if status != http.StatusRequestTimeout {
								http.Error(w, e.Error(), int(status))
							} else {
								writeMap(w, map[string]string{
									"token": t.String(),
								})
							}
							return
						},
					}
				}
			}
			return
		},
	}
}

func writeMap(w http.ResponseWriter, out map[string]string) error {

	// write out
	w.Header().Set("Content-Type", "application/json")
	js := json.NewEncoder(w)
	return js.Encode(out)
}
