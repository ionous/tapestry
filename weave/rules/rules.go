package rules

import (
	"fmt"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

// results of reading author specified pairs of pattern name and rule name.
// for example ("before someone jumping", "people jump for joy")
type RuleName struct {
	Pattern        []string // base pattern without any prefix or suffix
	Label          string   // friendly name of the rule itself
	Prefix         Prefix
	Suffix         Suffix
	ExcludesPlayer bool // true if the rule should apply to all actors
}

type RuleInfo struct {
	Name           string // name of the pattern / kind
	Label          string // friendly name of the rule itself
	Rank           int    // smaller ranked rules run first
	Stop           bool
	Jump           rt.Jump
	ExcludesPlayer bool // true if the rule should apply to all actors
}

func (n RuleName) IsEvent() bool {
	return slices.Contains(n.Pattern, kindsOf.Action.String())
}

// instead and report are grouped with before and after respectively
func (n RuleName) EventName() (ret string) {
	kind := n.Pattern[len(n.Pattern)-1]
	switch n.Prefix {
	case When:
		ret = kind
	case Instead:
		ret = event.BeforePhase.PatternName(kind)
	case Report:
		ret = event.AfterPhase.PatternName(kind)
	default:
		ret = n.Prefix.String() + " " + kind
	}
	return
}

// pattern name as specified
// optional rule name as specified
// ex. Define rule:named:do: ["activating", "the standard activating action" ]
func ReadPhrase(ks g.Kinds, patternSpec, ruleSpec string) (ret RuleName, err error) {
	patternSpec = inflect.Normalize(patternSpec)
	name, suffix := findSuffix(patternSpec)
	// return name sans any prefix, and any prefix the name had.
	short, prefix := findPrefix(name)
	var excludesPlayer bool
	const someone = "someone "
	if excludesPlayer = strings.HasPrefix(short, someone); excludesPlayer {
		short = short[len(someone):]
	}
	// fix: probably want some sort of "try prefix/suffix" that attempts to chop the parts
	// but restores them if it cant find them --
	// maybe see jess -- it does that sort of partial matching
	// and itd be neat to be able to use it here.
	// fix: we can pass in the base type
	if k, e := ks.GetKindByName(short); e != nil {
		err = e
	} else if pat := g.Ancestry(k); len(pat) == 0 {
		err = fmt.Errorf("couldnt determine ancestry of %q", short)
	} else {
		ret = RuleName{
			Pattern:        pat,
			Label:          inflect.Normalize(ruleSpec),
			Prefix:         prefix,
			Suffix:         suffix,
			ExcludesPlayer: excludesPlayer,
		}
	}
	return
}

// match an author specified pattern reference
// to various naming conventions and pattern definitions
// to determine the intended pattern name, rank, and termination behavior.
// for example: "instead of x", "before x", "after x", "report x".
func (n RuleName) GetRuleInfo() (ret RuleInfo, err error) {
	kind := n.Pattern[len(n.Pattern)-1]
	switch pattern.Categorize(n.Pattern) {
	default:
		err = errutil.Fmt("can't have a %q event", kind)

	// for regular patterns, supports sorting rules before/after
	case pattern.Calls:
		switch n.Prefix {
		case Instead, Report:
			err = errutil.Fmt("%q isn't an action and doesn't support %q", kind, n.Prefix)
		default:
			// ex. "before normal pattern", return "normal pattern"
			ret = RuleInfo{Name: kind, Rank: n.Prefix.rank()}
		}

	case pattern.Sends:
		stop, jump := n.Prefix.stopJump()
		ret = RuleInfo{
			Name:           n.EventName(),
			Rank:           n.Prefix.rank(),
			Stop:           stop,
			Jump:           jump,
			ExcludesPlayer: n.ExcludesPlayer,
		}
	}
	// suffix will override any stop/jump settings
	if err == nil {
		switch n.Suffix {
		case Jumps:
			ret.Stop = false
			ret.Jump = rt.JumpNow
		case Stops:
			ret.Stop = true
			ret.Jump = rt.JumpNow
		case Continues:
			ret.Stop = false
			ret.Jump = rt.JumpLater
		}
		ret.Label = n.Label
	}
	return
}

// Check to see if there are counters that might need updating on the regular.
// tdb: could this be processed at load time (storyImport)
func DoesUpdate(exe []rt.Execute) (okay bool) {
	// wrap the exes up into a single block for easier searching
	var m rtti.Execute_Slots = exe
	return searchCounters(&m)
}

// tdb: could this? be processed at load time (storyImport)
func DoesTerminate(exe []rt.Execute) bool {
	var terminal bool //provisionally continues
Out:
	for _, el := range exe {
		switch el := el.(type) {
		case *debug.DebugLog:
			// skip comments and debug logs
			// todo: make a "no op" interface so other things can join in?
			continue

		case *core.ChooseBranch:
			for next := el.Else; next != nil; {
				switch b := next.(type) {
				case *core.ChooseBranch:
					next = b.Else
				case *core.ChooseNothingElse:
					terminal = true
					break Out
				default:
					panic(errutil.Sprintf("unknown type of branch %T", next))
				}
			}

		default:
			// any statement other than a log or branch terminates
			terminal = true
			break Out
		}
	}
	return terminal
}

// return the first response definition in the block
func FindNamedResponse(exe []rt.Execute) (ret string) {
	var m rtti.Execute_Slots = exe
	if op, e := searchForFlow(&m, &render.Zt_RenderResponse); e != nil {
		panic(e)
	} else if response, ok := op.(*render.RenderResponse); ok && response.Text != nil {
		ret = response.Name
	}
	return
}
