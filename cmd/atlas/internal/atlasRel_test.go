package internal

import (
	"os"

	"git.sr.ht/~ionous/tapestry/tables"
)

func ExampleRelData() {
	Templates.ExecuteTemplate(os.Stdout, "relList", []Relation{
		{"containing", "containers", tables.ONE_TO_MANY, "things", "Containers contain stuff."},
		{"driving", "people", tables.ONE_TO_ONE, "cars", "No backseat drivers please."},
	})

	// Output:
	// <h1>Relations</h1>
	// <dl>
	//   <dt><a href="/atlas/relations/containing">Containing</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#containers">Containers</a> to many <a href="/atlas/kinds#things">Things</a>.
	//  Containers contain stuff.</dd>
	//   <dt><a href="/atlas/relations/driving">Driving</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#people">People</a> to <a href="/atlas/kinds#cars">Cars</a>.
	//  No backseat drivers please.</dd>
	// </dl>
}

// FIX: database format changed
func xExampleRelDB() {
	db := tables.CreateTest("ExampleRelDB", true)
	defer db.Close()
	if e := CreateTestData(db); e != nil {
		panic(e)
	} else if e := CreateAtlas(db); e != nil {
		panic(e)
	} else if e := listOfRelations(os.Stdout, db); e != nil {
		panic(e)
	}

	// Output:
	// <h1>Relations</h1>
	// <dl>
	//   <dt><a href="/atlas/relations/containing">Containing</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#vehicles">Vehicles</a> to many <a href="/atlas/kinds#people">People</a>.
	//  The outside of insides.</dd>
	// </dl>
}
