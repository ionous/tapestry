package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	. "git.sr.ht/~ionous/iffy/cmd/atlas/internal"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/web"
	"git.sr.ht/~ionous/iffy/web/support"
	_ "github.com/mattn/go-sqlite3"
)

// go run atlas.go -in /Users/ionous/Documents/Iffy/scratch/shared/play.db
func main() {
	var fileName string
	flag.StringVar(&fileName, "in", "", "input file name (sqlite3)")
	flag.Parse()
	if len(fileName) == 0 || fileName == "memory" {
		fileName = "file:test.db?cache=shared&mode=memory"
	}
	if db, e := sql.Open(tables.DefaultDriver, fileName); e != nil {
		log.Fatalln("db open", e)
	} else if _ /*fix temp view*/ = CreateAtlas(db); e != nil {
		log.Fatalln("db view", e)
	} else {
		if fileName == "memory" {
			if e := CreateTestData(db); e != nil {
				log.Fatal(e)
			}
		}

		m := http.NewServeMux()
		m.HandleFunc("/atlas/", web.HandleResource(Atlas(db)))
		go support.OpenBrowser("http://localhost:8080/atlas/")
		log.Fatal(http.ListenAndServe(":8080", m))
	}
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
