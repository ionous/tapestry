package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// "an {activity:patterned_activity} or a {value:variable_type}");
func imp_pattern_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	err = reader.Option(r, "pattern_type", reader.ReadMaps{
		"$ACTIVITY": func(m reader.Map) (err error) {
			ret, err = imp_patterned_activity(k, m)
			return
		},
		"$VALUE": func(m reader.Map) (err error) {
			ret, err = imp_variable_type(k, m)
			return
		},
	})
	return
}

// make.str("patterned_activity", "{activity}");
// returns "prog" as the name of a type  ( eases the difference b/t user named kinds, and internally named types )
func imp_patterned_activity(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if e := reader.Const(r, "patterned_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = k.NewName("execute", tables.NAMED_TYPE, reader.At(r))
	}
	return
}

// make.run("pattern_actions", ... "To {name:pattern_name}: {+pattern_rule}")
func imp_pattern_actions(k *Importer, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_actions"); e != nil {
		err = e
	} else if patternName, e := imp_pattern_name(k, op.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = imp_pattern_rules(k, patternName, op.MapOf("$PATTERN_RULES"))
	}
	return
}

func imp_pattern_rules(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_rules"); e != nil {
		err = e
	} else {
		err = reader.Repeats(op.SliceOf("$PATTERN_RULE"),
			func(el reader.Map) (err error) {
				return imp_pattern_rule(k, patternName, el)
			})
	}
	return
}

func imp_pattern_rule(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_rule"); e != nil {
		err = e
	} else if rule, e := imp_pattern_hook(k, op.MapOf("$HOOK")); e != nil {
		err = e
	} else if i, e := k.DecodeSlot(op.MapOf("$GUARD"), "bool_eval"); e != nil {
		err = e
	} else {
		rule.addFilter(i.(rt.BoolEval))
		if patternProg, e := k.NewProg(rule.typeName(), rule.buildRule()); e != nil {
			err = e
		} else {
			k.NewPatternRule(patternName, patternProg)
		}
	}
	return
}

// opt("pattern_hook", "{activity} or {result:pattern_return}");
func imp_pattern_hook(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "pattern_hook", reader.ReadMaps{
		"$ACTIVITY": func(m reader.Map) (err error) {
			if act, e := imp_activity(k, m); e != nil {
				err = e
			} else {
				ret = newExecuteRule(act)
			}
			return
		},
		"$RESULT": func(m reader.Map) (err error) {
			ret, err = imp_pattern_return(k, m)
			return
		},
	})
	return
}

// run("pattern_return", "return {result:pattern_result}");
// note: this slat exists for composer formatting reasons only...
func imp_pattern_return(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := reader.Unpack(r, "pattern_return"); e != nil {
		err = e
	} else {
		ret, err = imp_pattern_result(k, m.MapOf("$RESULT"))
	}
	return
}

// opt("pattern_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
func imp_pattern_result(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "pattern_result", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			ret, err = imp_primitive_func(k, m)
			return
		},
		"$OBJECT": func(m reader.Map) (err error) {
			ret, err = imp_object_func(k, m)
			return
		},
	})
	return
}

// opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}")
func imp_primitive_func(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "primitive_func", reader.ReadMaps{
		"$NUMBER_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "number_eval"); e != nil {
				err = e
			} else {
				ret = newNumberRule(i.(rt.NumberEval))
			}
			return
		},
		"$TEXT_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "text_eval"); e != nil {
				err = e
			} else {
				ret = newTextRule(i.(rt.TextEval))
			}
			return
		},
		"$BOOL_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "bool_eval"); e != nil {
				err = e
			} else {
				ret = newBoolRule(i.(rt.BoolEval))
			}
			return
		},
	})
	return
}

// run("object_func", "an object named {name%text_eval}");
func imp_object_func(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := reader.Unpack(r, "object_func"); e != nil {
		err = e
	} else if i, e := k.DecodeSlot(m, "text_eval"); e != nil {
		err = e
	} else {
		// FIX: we should wrap the text with a runtime "test the object matches" command
		// and, -- for simple text values -- add ephemera for "name, type"
		ret = newTextRule(i.(rt.TextEval))
	}
	return
}
