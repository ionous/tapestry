package internal

import (
	"database/sql"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// import an object type description
func TestImpObjectType(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if n, e := imp_object_type(k, _object_type); e != nil {
		t.Fatal(e)
	} else if n.String() != "animals" {
		t.Fatal(n)
	}
}

// import a variable type description
func TestImpVariableTypePrimitive(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if varType, e := imp_variable_type(k, _variable_type); e != nil {
		t.Fatal(e)
	} else if varType.String() != tables.EVAL_EXPR {
		t.Fatal(varType)
	}
}

// import a variable declaration
func TestImpVariableDeclObject(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if varName, typeName, e := imp_variable_decl(k, _variable_decl); e != nil {
		t.Fatal(e)
	} else if varName.String() != "pet" {
		t.Fatal(varName)
	} else if typeName.String() != "animals" {
		t.Fatal(typeName)
	}
}

func TestImpPatternVariablesDecl(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_variables_decl(k, _pattern_variables_decl); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 3)
		if diff := pretty.Diff(buf.String(), lines(
			"corral,pattern_name",  // 1
			"pet,variable_name",    // 2
			"animals,plural_kinds", // 3
			"1,2,3",
		)); len(diff) > 0 {
			t.Fatal("mismatch", diff)
		} else {
			t.Log("ok")
		}
	}
}

func TestImpPrimitiveType(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_primitive_type(k, _primitive_type); e != nil {
		t.Fatal(e)
	} else if typ.String() != "bool" {
		t.Fatal(typ)
	}
}

func TestImpPatternType_Activity(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_pattern_type(k, _pattern_type_activity); e != nil {
		t.Fatal(e)
	} else if typ.String() != "prog" {
		t.Fatal(typ)
	}
}

func TestImpPatternType_Primitive(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_pattern_type(k, _pattern_type_primitive); e != nil {
		t.Fatal(e)
	} else if typ.String() != "bool" {
		t.Fatal(typ)
	}
}

func TestImpPatternName(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if n, e := imp_pattern_name(k, _pattern_name); e != nil {
		t.Fatal(e)
	} else if n.String() != "corral" {
		t.Fatal(n)
	}
}

func TestImpPatternDecl(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_decl(k, _pattern_decl); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 3)
		if diff := pretty.Diff(buf.String(), lines(
			"corral,pattern_name", // 1
			"prog,type",           // 2
			"1,1,2",
		)); len(diff) > 0 {
			t.Fatal("mismatch", diff)
		} else {
			t.Log("ok")
		}
	}
}

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func newTestImporter(t *testing.T) (ret *Importer, retDB *sql.DB) {
	const path = "file:test.db?cache=shared&mode=memory"
	// if path, e := getPath(t.Name() + ".db"); e != nil {
	// 	t.Fatal(e)
	// } else
	if db, e := sql.Open("sqlite3", path); e != nil {
		t.Fatal("db open", e)
	} else {
		if e := tables.CreateEphemera(db); e != nil {
			t.Fatal("create ephemera", e)
		} else {
			ret, retDB = NewImporter(t.Name(), db), db
		}
	}
	return
}

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

var _pattern_decl = map[string]interface{}{
	"id":   "id-171a4d90ca7-5",
	"type": "pattern_decl",
	"value": map[string]interface{}{
		"$NAME": _pattern_name,
		"$TYPE": _pattern_type_activity,
	},
}

var _pattern_name = map[string]interface{}{
	"id":    "id-171a4d90ca7-3",
	"type":  "pattern_name",
	"value": "corral",
}

var _pattern_type_activity = map[string]interface{}{
	"id":   "id-171a4d90ca7-4",
	"type": "pattern_type",
	"value": map[string]interface{}{
		"$ACTIVITY": map[string]interface{}{
			"id":    "id-171a4d90ca7-6",
			"type":  "pattern_activity",
			"value": "$ACTIVITY",
		},
	},
}

var _pattern_type_primitive = map[string]interface{}{
	"id":   "id-171a8ba0566-4",
	"type": "pattern_type",
	"value": map[string]interface{}{
		"$VALUE": map[string]interface{}{
			"id":   "id-171a8ba0566-6",
			"type": "variable_type",
			"value": map[string]interface{}{
				"$PRIMITIVE": _primitive_type,
			},
		},
	},
}

var _primitive_type = map[string]interface{}{
	"id":    "id-171a8ba0566-7",
	"type":  "primitive_type",
	"value": "$BOOL",
}

var _pattern_variables_decl = map[string]interface{}{
	"id":   "id-1719a47c939-7",
	"type": "pattern_variables_decl",
	"value": map[string]interface{}{
		"$PATTERN_NAME": map[string]interface{}{
			"id":    "id-1719a47c939-3",
			"type":  "pattern_name",
			"value": "corral",
		},
		"$VARIABLE_DECL": []interface{}{
			_variable_decl,
		},
	},
}

var _variable_decl = map[string]interface{}{
	"id":   "id-1719a47c939-11",
	"type": "variable_decl",
	"value": map[string]interface{}{
		"$TYPE": map[string]interface{}{
			"id":   "id-1719a47c939-9",
			"type": "variable_type",
			"value": map[string]interface{}{
				"$OBJECT": _object_type,
			},
		},
		"$NAME": map[string]interface{}{
			"id":    "id-1719a47c939-10",
			"type":  "variable_name",
			"value": "pet",
		},
	},
}

var _variable_type = map[string]interface{}{
	"id":   "id-1719a47c939-4",
	"type": "variable_type",
	"value": map[string]interface{}{
		"$PRIMITIVE": map[string]interface{}{
			"id":    "id-1719a47c939-8",
			"type":  "primitive_type",
			"value": "$TEXT",
		},
	},
}

var _object_type = map[string]interface{}{
	"id":   "id-1719a47c939-14",
	"type": "object_type",
	"value": map[string]interface{}{
		"$AN": map[string]interface{}{
			"id":    "id-1719a47c939-12",
			"type":  "an",
			"value": "$AN",
		},
		"$KINDS": map[string]interface{}{
			"id":    "id-1719a47c939-13",
			"type":  "plural_kinds",
			"value": "animals",
		},
	},
}
