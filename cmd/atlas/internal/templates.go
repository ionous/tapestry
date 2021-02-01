package internal

import (
	"database/sql"
	"html/template"
	"reflect"
	"regexp"
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

var AtlasTemplate = atlasTemplate
var Templates *template.Template = template.New("none").Funcs(funcMap)

//go:generate templify -p main -o atlas.gen.go atlas.sql
func CreateAtlas(db *sql.DB) (err error) {
	if _, e := db.Exec(AtlasTemplate()); e != nil {
		err = errutil.New("CreateAtlas:", e)
	}
	return
}

var spaces = regexp.MustCompile(`\s+`)

var funcMap = template.FuncMap{
	"title": lang.Titlecase,
	"safe": func(s string) string {
		return spaces.ReplaceAllString(s, "-")
	},
	"prefix": strings.HasPrefix,
	"suffix": strings.HasSuffix,
	// return true if the struct field in els before idx differs from the one at idx
	"changing": func(idx int, field string, els reflect.Value) (ret bool) {
		if idx == 0 {
			ret = true
		} else {
			curr, prev := els.Index(idx), els.Index(idx-1)
			c := curr.Elem().FieldByName(field).Interface()
			p := prev.Elem().FieldByName(field).Interface()
			ret = c != p
		}
		return
	},
}

func registerTemplate(n, t string) {
	Templates = template.Must(Templates.New(n).Parse(t))
}
