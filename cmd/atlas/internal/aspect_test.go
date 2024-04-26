package internal

import (
	"os"

	"git.sr.ht/~ionous/tapestry/tables"
)

func ExampleAspectData() {
	Templates.ExecuteTemplate(os.Stdout, "aspectList", []Aspect{{
		Name:  "flightiness",
		Kinds: []string{"vehicles", "birds"},
		Traits: []Trait{{
			Name: "flightless",
			Spec: "grounded.",
		}, {
			Name: "glide worthy",
			Spec: "Better at landing than taking off.",
		}, {
			Name: "flight worthy",
		}},
	}, {
		Name:  "bool",
		Spec:  "example binary aspect.",
		Kinds: []string{"things"},
		Traits: []Trait{{
			Name: "true",
		}, {
			Name: "false",
		}},
	}})

	// Output:
	// <h1>Aspects</h1>
	//   <a href="#flightiness">Flightiness</a>,
	//   <a href="#bool">Bool</a>.
	// <h2 id="flightiness">Flightiness</h2>
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#vehicles">Vehicles</a>,
	//   <a href="/atlas/kinds#birds">Birds</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>Flightless</dt>
	//    <dd>grounded.</dd>
	//   <dt>Glide Worthy</dt>
	//    <dd>Better at landing than taking off.</dd>
	//   <dt>Flight Worthy</dt>
	// </dl>
	// <h2 id="bool">Bool</h2>
	// example binary aspect.
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#things">Things</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>True</dt>
	//   <dt>False</dt>
	// </dl>
}

// FIX: database format changed
func xExampleAspectDB() {
	db := tables.CreateTest("ExampleAspectDB", true)
	defer db.Close()
	if e := CreateTestData(db); e != nil {
		panic(e)
	} else if e := CreateAtlas(db); e != nil {
		panic(e)
	} else if e := listOfAspects(os.Stdout, db); e != nil {
		panic(e)
	}

	// Output:
	// <h1>Aspects</h1>
	//   <a href="#flightiness">Flightiness</a>.
	// <h2 id="flightiness">Flightiness</h2>
	// The flight worthiness of vehicles, an example of an aspect with several traits.
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#vehicles">Vehicles</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>Flight Worthy</dt>
	//   <dt>Flightless</dt>
	//   <dt>Glide Worthy</dt>
	//    <dd>Better at landing than taking off.</dd>
	// </dl>
}
