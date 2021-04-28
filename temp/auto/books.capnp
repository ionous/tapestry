using Go = import "/go.capnp";
using Opt = import "options.capnp";
$Go.package("books");
$Go.import("foo/books");


struct Circle $baz("class") {
	title @0 :Text $baz("field");
}


struct Square $Go.name("Squeak") {
	width @0 : Int32;
}

struct Book {
	title @0 :Text;
	# Title of the book.

	pageCount @1 :Int32;
	# Number of pages in the book.

	union {
	    circle @2 :Circle;      # radius
	    square @3 :Square;      # width
	  }
}