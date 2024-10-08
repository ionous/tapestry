package internal

import (
	"database/sql"
	"fmt"
	"io"

	"git.sr.ht/~ionous/tapestry/tables"
)

type Pairing struct {
	Rel   *Relation
	Pairs []*Pair
}

type Pair struct {
	First, Second string
}

func listOfPairs(w io.Writer, relation string, db *sql.DB) (err error) {
	var rel Relation
	var pair Pair
	var pairs []*Pair

	if e := tables.QueryAll(db,
		fmt.Sprintf(`
		select relation, kind, cardinality, otherKind, coalesce((
			select spec from mdl_spec 
			where type='relation' and name=relation
			limit 1), '')
		from mdl_rel
		where relation = '%s'`, relation),
		func() (err error) {
			return
		}, &rel.Name, &rel.Kind, &rel.Cardinality, &rel.OtherKind, &rel.Spec); e != nil {
		err = e
	} else if e := tables.QueryAll(db,
		fmt.Sprintf(`
		select noun, otherNoun
		from mdl_pair
		where relation = '%s'`, relation),
		func() (err error) {
			pin := pair
			pairs = append(pairs, &pin)
			return
		}, &pair.First, &pair.Second,
	); e != nil {
		err = e
	} else {
		pin := rel
		err = Templates.ExecuteTemplate(w, "pairList", &Pairing{
			Rel:   &pin,
			Pairs: pairs,
		})
	}
	return
}

func init() {
	registerTemplate("relHeader", `Relates
	{{- if prefix .Cardinality "any_" }} many
	{{- end -}}
{{- "" }} <a href="/atlas/kinds#{{.Kind|safe}}">{{.Kind|title}}</a> to
	{{- if suffix .Cardinality "_any" }} many
	{{- end -}} 
{{- "" }} <a href="/atlas/kinds#{{.OtherKind|safe}}">{{.OtherKind|title}}</a>.
	{{- if .Spec }}
{{ "" }} {{ .Spec }}
	{{- end -}}
`)

	registerTemplate("pairList", `
<h1>{{.Rel.Name|title}}</h1>
{{ template "relHeader" .Rel }}
<table>
	{{- range $i, $el := .Pairs }}
<tr>
  <td>{{ if changing $i "First" $.Pairs }}<a href="/atlas/nouns#{{.First|safe}}">{{.First|title}}</a>{{end}}</td>
  <td>{{ if changing $i "Second" $.Pairs }}<a href="/atlas/nouns#{{.Second|safe}}">{{.Second|title}}</a>{{end}}</td>
</tr>
	{{- end }}
</table>
`)
}
