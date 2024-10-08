#ifndef TAPESETRY_H
#define TAPESETRY_H

#include <godot_cpp/classes/sprite2d.hpp>

namespace godot {

class Tapestry : public Node {
	GDCLASS(Tapestry, Node)

public:
	static Variant post( const String& endpoint, const String &json );

protected:
	static void _bind_methods();
};

}

#endif