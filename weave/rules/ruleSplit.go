package rules

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"golang.org/x/exp/slices"
)

// chop a series of statements into separate rules
// if (and only if) the statements are a set of if-else statements
// returns true if any of the rules updates
func Split(name string, rank int, exe []rt.Execute) (rules []mdl.Rule) {
	if tree := findTree(exe); tree != nil {
		rules, _ = splitTree(name, rank, tree)
	} else {
		// no tree? then one terminating rule.
		rules = append(rules, mdl.Rule{
			Name: name,
			Rank: rank,
			Prog: assign.Prog{
				Exe:        exe,
				Terminates: true,
			},
		})
	}
	return
}

// first scan for a branching statement tree
func findTree(exe []rt.Execute) (ret core.Brancher) {
FindTree:
	for _, el := range exe {
		switch el := el.(type) {
		case *debug.DebugLog:
			// skip debug logs when trying to find a tree
		case *core.ChooseBranch:
			// a branch; we'll use it if we can
			ret = el
		default:
			// some other statement
			ret = nil
			break FindTree
		}
	}
	return
}

// split all the parts of the passed tree into separate rules
func splitTree(name string, rank int, tree core.Brancher) ([]mdl.Rule, bool) {
	var rules []mdl.Rule
	var update updateTracker
	for i, next := 1, tree; next != nil; i++ {
		ruleName := name
		if len(name) > 0 {
			ruleName = fmt.Sprintf("%s (%d)", name, i)
		}
		switch el := next.(type) {
		case *core.ChooseBranch:
			next = el.Else // a continuable rule
			rules = append(rules, mdl.Rule{
				Name: ruleName,
				Rank: rank,
				Prog: assign.Prog{
					Filter: el.If,
					Args:   el.Args,
					Exe:    el.Exe,
					Updates: update.CheckFilter(el.If) ||
						update.CheckArgs(el.Args),
				}})

		case *core.ChooseNothingElse:
			next = nil // terminal rule
			rules = append(rules, mdl.Rule{
				Name: ruleName,
				Rank: rank,
				Prog: assign.Prog{
					Filter:     nil, //
					Args:       el.Args,
					Exe:        el.Exe,
					Updates:    update.CheckArgs(el.Args),
					Terminates: true,
				},
			})
		} // ~ end switch
	} // ~ end for
	// rules are read from the db in reverse order.
	slices.Reverse(rules)
	return rules, update.HasUpdate()
}
