/*
 * CODE GENERATED AUTOMATICALLY WITH
 *    github.com/wlbr/templify
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 */

package tables

// runTemplate is a generated function returning the template as a string.
// That string should be parsed by the functions of the golang's template package.
func runTemplate() string {
	var tmpl = "/**\n" +
		" * for saving, restoring a player's game session.\n" +
		" */\n" +
		"create table if not exists \n" +
		"\trun_domain( domain int, active int, primary key( domain )); \n" +
		"\n" +
		"\n" +
		"/* we dont need an \"active\" -- we can join against run_domain, or write 0 to domain to disable a pair. */\n" +
		"create table if not exists \n" +
		"\trun_pair( domain int, relKind int, oneNoun int, otherNoun int, unique( relKind, oneNoun, otherNoun ) ); \n" +
		"\n" +
		"/**\n" +
		" * find the id and name of all active domains\n" +
		" * path isnt really needed because any parts of an activated domain are themselves individually active.\n" +
		" */\n" +
		"create view if not exists\n" +
		"domain_scope as \n" +
		"select md.rowid as domain, md.domain as name\n" +
		"from run_domain rd \n" +
		"join mdl_domain md \n" +
		"\ton rd.active > 0 and rd.domain = md.rowid;\n" +
		"\n" +
		"create view if not exists\n" +
		"kind_scope as \n" +
		"select mk.rowid as kind, mk.kind as name, mk.path\n" +
		"from domain_scope\n" +
		"join mdl_kind mk \n" +
		"\tusing (domain);\n" +
		"\n" +
		"create view if not exists\n" +
		"noun_scope as \n" +
		"select mn.rowid as noun, mn.noun as name, mn.kind\n" +
		"from domain_scope\n" +
		"join mdl_noun mn \n" +
		"\tusing (domain);\n" +
		"\n" +
		"/* returns the kind's id as 'relKind' */ \n" +
		"create view if not exists\n" +
		"rel_scope as \n" +
		"select mk.rowid as relKind, mk.kind as name, mr.oneKind, mr.otherKind, mr.cardinality\n" +
		"from domain_scope\n" +
		"join mdl_kind mk \n" +
		"\tusing (domain)\n" +
		"join mdl_rel mr\n" +
		"\ton (mk.rowid = mr.relKind);\n" +
		"\n" +
		"\n" +
		"/* */\n" +
		"create view if not exists\n" +
		"rp_names as\n" +
		"select rp.domain, mk.kind as relName, rp.relKind, one.noun as oneName, rp.oneNoun, other.noun as otherName, rp.otherNoun, rel.cardinality\n" +
		"from run_pair rp\n" +
		"join mdl_kind mk\n" +
		"\ton (mk.rowid = rp.relKind)\n" +
		"join mdl_noun one\n" +
		"\ton (one.rowid = rp.oneNoun)\n" +
		"join mdl_noun other\n" +
		"\ton (other.rowid = rp.otherNoun)\n" +
		"join mdl_rel rel \n" +
		"\ton (rel.relKind = rp.relKind)\n" +
		"where rp.domain > 0;"
	return tmpl
}
