package mdl

import "git.sr.ht/~ionous/tapestry/affine"

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler interface {
  // the pattern half of Start; domain, kind, field are a pointer into Field
  // value should be a marshaled compact value
  // fix? [ move marshaling into the implementation? ]
  Assign(domain, kind, field, value string) (err error)
  // author tests of stories
  Check(domain, name, value string, affinity affine.Affinity, prog, at string) (err error)
  // pairs of domain name and (domain) dependencies
  Domain(domain, requires, at string) (err error)
  // note: the domain exists to uniquely identify the kind;
  // it's not actually stored in the field table and requires the write to transform it properly.
  Field(domain, kind, field string, affinity affine.Affinity, typeName, at string) (err error)
  Grammar(domain, name, prog, at string) (err error)
  // singular name of kind and materialized hierarchy of ancestors separated by commas
  Kind(domain, kind, path, at string) (err error)
  // words for authors and game players refer to nouns
  // follows the domain rules of Noun.
  Name(domain, noun, name string, rank int, at string) (err error)
  // the domain tells the scope in which the noun was defined
  // ( the same as - or a child of - the domain of the kind ) (err error)
  Noun(domain, noun, kind, at string) (err error)
  //
  Opposite(domain, oneWord, otherWord, at string) (err error)
  // domain captures the scope in which the pairing was defined.
  // within that scope: the noun, relation, and otherNoun are all unique names --
  // even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
  Pair(domain, relKind, oneNoun, otherNoun, at string) (err error)
  // doesn't store at because its kind already defines that
  Pat(domain, kind, labels, result string) (err error)
  // a plural word (many) can have at most one singular definition per domain
  // ie. people and persons are valid plurals of person,
  // but people as a singular can only be defined as person ( not also human ) (err error)
  Plural(domain, many, one, at string) (err error)
  // relation and constraint between two kinds of nouns
  //  fix? the data is duplicated in kinds and fields... should this be removed?
  // might also consider adding a cardinality field to the relation kind, and then use init for individual relations
  Rel(domain, relKind, oneKind, otherKind, cardinality, at string) (err error)
  //
  Rule(domain, pattern, target string, phase int, filter, prog, at string) (err error)
  // the noun half of what was Start.
  // domain, noun, field reference a join of Noun and Kind to get a filtered Field.
  Value(domain, noun, field, value, at string) (err error)
}
