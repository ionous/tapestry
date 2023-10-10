#include "tapext.h"
#include <godot_cpp/core/class_db.hpp>
#include <godot_cpp/classes/json.hpp>
#include "taplib.h"

using namespace godot;

//---------------------------------------------------------------------------
void Tapestry::_bind_methods() {
	ClassDB::bind_static_method("Tapestry", D_METHOD("post", "endpoint", "json"), &Tapestry::post);
}

//---------------------------------------------------------------------------
Variant Tapestry::post( const String& endpoint, const String& json ) {
	// string literals in go are UTF-8 because the source files are defined as utf8.
	// json is technically utf8, all input would be in ansi. anyway, this is fine if slow.
	CharString jsChars = json.utf8();   // copies
	GoString jsGo = { jsChars.ptr(), jsChars.length() };
	//
	CharString endChars = endpoint.utf8(); // copies
	GoString endGo = { endChars.ptr(), endChars.length() };

	// go returns its own memory; pinned until the next time post is called.
	const char * got = Post(endGo, jsGo);
	// String() copies, and parse_string() convert to arrays and dictionaries so the memory is good to go on return anyway.
	String msg( got );
	//
	Variant res = JSON::parse_string( msg );
	if (!res) {
		WARN_PRINT_ONCE_ED(got);
	}
	return res;
}
