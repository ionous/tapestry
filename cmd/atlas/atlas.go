package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	. "git.sr.ht/~ionous/tapestry/cmd/atlas/internal"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web"
	"git.sr.ht/~ionous/tapestry/web/support"
	_ "github.com/mattn/go-sqlite3"
)

// go run atlas.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	var fileName string
	flag.StringVar(&fileName, "in", "", "input file name (sqlite3)")
	flag.Parse()
	if db, e := openDB(fileName); e != nil {
		panic(e)
	} else {
		defer db.Close()
		if _ /*fix temp view*/ = CreateAtlas(db); e != nil {
			log.Fatalln("db view", e)
		} else {
			m := http.NewServeMux()
			m.HandleFunc("/atlas/", web.HandleResource(Atlas(db)))
			go support.OpenBrowser("http://localhost:8080/atlas/")
			log.Fatal(http.ListenAndServe(":8080", m))
		}
	}
}

func openDB(name string) (ret *sql.DB, err error) {
	useTestData := len(name) == 0 || name == "memory"
	if !useTestData {
		ret, err = tables.CreateRunTime(name)
	} else {
		db := tables.CreateTest("testdata", true)
		if e := CreateTestData(db); e != nil {
			db.Close()
			err = e
		} else {
			ret = db
		}
	}
	return
}

func Atlas(db *sql.DB) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "atlas":
				ret = &web.Wrapper{
					Finds: func(name string) (ret web.Resource) {
						switch name {
						case "nouns":
							ret = Empty(name)
						case "kinds":
							return &web.Wrapper{
								Gets: func(ctx context.Context, w http.ResponseWriter) error {
									return ListOfKinds(w, db)
								},
							}
						}
						return
					},
					Gets: func(ctx context.Context, w http.ResponseWriter) error {
						return Templates.ExecuteTemplate(w, "links", []struct{ Link, Text string }{
							{"/atlas/kinds/", "kinds"},
							{"/atlas/nouns/", "nouns"},
						})
					},
				}
			}
			return
		},
	}
}

func Empty(name string) web.Resource {
	return &web.Wrapper{
		Gets: func(ctx context.Context, w http.ResponseWriter) error {
			_, e := fmt.Fprintf(w, "No %s to see here", name)
			return e
		},
	}
}
