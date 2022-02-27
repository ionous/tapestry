package chart

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"os"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func Chart(inFile, outFile, scope string) (ret int, err error) {
	if outf, e := os.Create(outFile); e != nil {
		err = e
	} else {
		defer outf.Close()
		if db, e := sql.Open(tables.DefaultDriver, inFile+"?mode=ro"); e != nil {
			err = errutil.New("couldn't create output file", inFile, e)
		} else {
			defer db.Close()
			var prep tables.Prep
			nounValue := prep.Prep(db, nounValueSql)
			nounName := prep.Prep(db, nounNameSql)
			if e := prep.Err(); e != nil {
				err = e
			} else {
				var rooms []Room
				var noun string
				var value []byte
				if rows, e := nounValue.Query(scope, "rooms", "compass", ""); e != nil {
					err = e
				} else if e := tables.ScanAll(rows, func() (err error) {
					var compass map[string]string
					if e := json.Unmarshal(value, &compass); e != nil {
						err = e
					} else {
						rooms = append(rooms, Room{noun, compass})
					}
					return
				}, &noun, &value); e != nil {
					err = e
				} else {
					funcMap := template.FuncMap{
						"title": strings.Title,
						"nameOf": func(noun string) (ret string, err error) {
							err = nounName.QueryRow(noun).Scan(&ret)
							return
						},
						"roomOfDoor": func(door string) (ret string, err error) {
							// note: js strings are quoted, but graphvis doesn't seem to care.
							if e := nounValue.QueryRow(scope, "doors", "destination", door).Scan(&door, &ret); e != nil {
								err = errutil.New(door, e)
							}
							return
						},
					}
					if t, e := template.New("chart").Funcs(funcMap).Parse(chartingTmpl); e != nil {
						err = e
					} else if e := t.Execute(outf, rooms); e != nil {
						err = e
					}
				}
			}
		}
	}
	return
}

type Room struct {
	Noun    string
	Compass map[string]string
}

// give a domain, kind, field:
// return all nouns of that kind, and the value for each noun of that field.
//go:embed nounValue.sql
var nounValueSql string

//go:embed nounName.sql
var nounNameSql string

//go:embed charting.tmpl
var chartingTmpl string
