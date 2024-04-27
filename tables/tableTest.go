package tables

import (
	"database/sql"
	"github.com/ionous/errutil"
	"os"
	"os/user"
	"path"
	"strings"
)

// add model and runtime to the same db
// if you run the test as "go test ... -args write"
// it'll write the db out in your user (home) directory
func CreateTest(name string, includeRunTables bool) *sql.DB {
	db, e := createTestDB(name, includeRunTables)
	if e != nil {
		panic(e)
	}
	return db
}

func createTestDB(name string, includeRunTables bool) (ret *sql.DB, err error) {
	var driver string
	if includeRunTables {
		driver = tapestryDriver
	} else {
		driver = defaultDriver
	}
	if fileName, e := resolveTestFile(name); e != nil {
		err = e
	} else if db, e := open(driver, fileName); e != nil {
		err = e
	} else {
		if e := createTables(db, "model", "modelView"); e != nil {
			err = e
		} else {
			if includeRunTables {
				err = createTables(db, "runTables", "runViews")
			}
		}
		if err != nil {
			db.Close()
		} else {
			ret = db
		}
	}
	return
}

// if you run the test as "go test ... -args write"
// it'll write the db out in your user (home) directory
func resolveTestFile(name string) (ret string, err error) {
	if os.Args[len(os.Args)-1] != "write" {
		ret = memory
	} else {
		// tests can have slash in their name
		rest := strings.Replace(name, "/", ".", -1) + ".db"
		if user, e := user.Current(); e != nil {
			err = errutil.New(e, "for", name)
		} else {
			ret = path.Join(user.HomeDir, rest)
		}
	}
	return
}
