package mdl

import (
  "git.sr.ht/~ionous/tapestry/affine"
  "git.sr.ht/~ionous/tapestry/dl/assign"
  "git.sr.ht/~ionous/tapestry/dl/grammar"
  "git.sr.ht/~ionous/tapestry/dl/literal"
  "git.sr.ht/~ionous/tapestry/rt"
)

// Modeler wraps writing to the model table
// so the implementation can handle verifying dependent names when needed.
type Modeler interface {
  Aspect(domain, aspect, at string, traits []string) error
  // author tests of stories
  Check(domain, name string, value literal.LiteralValue, exe []rt.Execute, at string) error
  // the pattern half of Start; domain, kind, field are a pointer into Field
  // value should be a marshaled compact value
  Default(domain, kind, field string, value assign.Assignment) error
  // pairs of domain name and (domain) dependencies
  Domain(domain, requires, at string) error
  // note: the domain exists to uniquely identify the kind;
  // it's not actually stored in the field table and requires the write to transform it properly.
  Field(domain, kind, field string, affinity affine.Affinity, typeName, at string) error
  Grammar(domain, name string, d *grammar.Directive, at string) error
  // singular name of kind and materialized hierarchy of ancestors separated by commas
  Kind(domain, kind, path, at string) error
  // words for authors and game players refer to nouns
  // follows the domain rules of Noun.
  Name(domain, noun, name string, rank int, at string) error
  // the domain tells the scope in which the noun was defined
  // ( the same as - or a child of - the domain of the kind ) error
  Noun(domain, noun, kind, at string) error
  //
  Opposite(domain, oneWord, otherWord, at string) error
  // domain captures the scope in which the pairing was defined.
  // within that scope: the noun, relation, and otherNoun are all unique names --
  // even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
  Pair(domain, relKind, oneNoun, otherNoun, at string) error
  // doesn't store at because its kind already defines that
  Pat(domain, kind, labels, result string) error
  // a plural word (many) can have at most one singular definition per domain
  // ie. people and persons are valid plurals of person,
  // but people as a singular can only be defined as person ( not also human ) error
  Plural(domain, many, one, at string) error
  // relation and constraint between two kinds of nouns
  //  fix? the data is duplicated in kinds and fields... should this be removed?
  // might also consider adding a cardinality field to the relation kind, and then use init for individual relations
  Rel(domain, relKind, oneKind, otherKind, cardinality, at string) error
  //
  Rule(domain, pattern, target string, phase int, filter rt.BoolEval, exe []rt.Execute, at string) error
  // the noun half of what was Start.
  // domain, noun, field reference a join of Noun and Kind to get a filtered Field.
  Value(domain, noun, field string, value literal.LiteralValue, at string) error

  //
  FindCompatibleField(domain, noun, field string, aff affine.Affinity) (retName string, retClass string, err error)
}
