#include "tapext.h"
#include <godot_cpp/core/class_db.hpp>
#include <godot_cpp/classes/json.hpp>
#include "taplib.h"

using namespace godot;

//---------------------------------------------------------------------------
void Tapestry::_bind_methods() {
	ClassDB::bind_static_method("Tapestry", D_METHOD("post", "json_string"), &Tapestry::post);
}


//---------------------------------------------------------------------------
Variant Tapestry::post( const String& js ) {
	// string literals in go are UTF-8 because the source files are defined as utf8
	// json technically is... probably all input would be in ansi ... this is fine for now.
	// its all slow anyway.
	WARN_PRINT_ED("got a post");
	CharString charstr = js.utf8();
	GoString str = { charstr.ptr(), charstr.length() };
  WARN_PRINT_ED(charstr.ptr());
	const char * got = Post(str);
	WARN_PRINT_ED("returned some value");
	WARN_PRINT_ED(got ? got : "got was null");
	// godot::String copies and then JSON copies the sub bits and that's just the way it is.
	// ( but right now we are only getting raw text back and not json commands )
	// return Variant( JSON::parse_string( String( got ) ));
	String want = got;
	return Variant(want);
}
