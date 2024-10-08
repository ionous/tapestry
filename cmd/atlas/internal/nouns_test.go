package internal

import (
	"os"

	"git.sr.ht/~ionous/tapestry/tables"
)

func ExampleNounData() {
	Templates.ExecuteTemplate(os.Stdout, "nounList", []Noun{
		{Name: "pants", Kind: "clothes", Spec: "An important start to every day, the putting on of.",
			Props: []Prop{
				// possibly links back to the declaring kind?
				{Name: "color", Value: "red"},
				{Name: "holiness level", Value: "5"},
			},
			Relations: []string{"containing", "stolen"},
		},
		{Name: "shoes", Kind: "clothes", Spec: "To be used after the pants, obviously."},
	})

	// Output:
	// <h1>Nouns</h1>
	//   <a href="#pants">Pants</a>,
	//   <a href="#shoes">Shoes</a>.
	//
	// <h2 id="pants">Pants</h2>
	// <span>Kind: <a href="/atlas/kinds#clothes">Clothes</a>.</span> <span class="spec">An important start to every day, the putting on of.</span>
	//
	// <h3>Properties</h3>
	// <ul>
	//   <li>Color: <span>red.</span></li>
	//   <li>Holiness Level: <span>5.</span></li>
	// </ul>
	//
	// <h3>Relations</h3>
	//   <a href="/atlas/relations/containing">Containing</a>,
	//   <a href="/atlas/relations/stolen">Stolen</a>.
	//
	// <h2 id="shoes">Shoes</h2>
	// <span>Kind: <a href="/atlas/kinds#clothes">Clothes</a>.</span> <span class="spec">To be used after the pants, obviously.</span>
	//
}

// FIX: database format changed
func xExampleNounDB() {
	db := tables.CreateTest("ExampleNounDB", true)
	defer db.Close()
	if e := CreateTestData(db); e != nil {
		panic(e)
	} else if e := CreateAtlas(db); e != nil {
		panic(e)
	} else if e := listOfNouns(os.Stdout, db); e != nil {
		panic(e)
	}

	// Output:
	//
	// <h1>Nouns</h1>
	//   <a href="#dune-buggy">Dune Buggy</a>,
	//   <a href="#enterprise">Enterprise</a>,
	//   <a href="#picard">Picard</a>,
	//   <a href="#riker">Riker</a>.
	//
	// <h2 id="dune buggy">Dune Buggy</h2>
	// <span>Kind: <a href="/atlas/kinds#cars">Cars</a>.</span> <span class="spec"></span>
	//
	// <h3>Properties</h3>
	// <ul>
	//   <li>Dune Buggy: <span>3.</span></li>
	// </ul>
	//
	// <h3>Relations</h3>
	//   <a href="/atlas/relations/containing">Containing</a>.
	//
	// <h2 id="enterprise">Enterprise</h2>
	// <span>Kind: <a href="/atlas/kinds#vehicles">Vehicles</a>.</span> <span class="spec"></span>
	//
	// <h3>Relations</h3>
	//   <a href="/atlas/relations/containing">Containing</a>.
	//
	// <h2 id="picard">Picard</h2>
	// <span>Kind: <a href="/atlas/kinds#people">People</a>.</span> <span class="spec"></span>
	//
	// <h3>Relations</h3>
	//   <a href="/atlas/relations/containing">Containing</a>.
	//
	// <h2 id="riker">Riker</h2>
	// <span>Kind: <a href="/atlas/kinds#people">People</a>.</span> <span class="spec"></span>
	//
	// <h3>Relations</h3>
	//   <a href="/atlas/relations/containing">Containing</a>.
	//
}
