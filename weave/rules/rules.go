package rules

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

// results of reading author specified pairs of pattern name and rule name.
// for example ("before someone jumping", "people jump for joy")
type RuleName struct {
	Short          string // base pattern without any prefix or suffix
	Label          string // friendly name of the rule itself
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

// instead and report are grouped with before and after respectively
func (n RuleName) EventName() (ret string) {
	switch n.Prefix {
	case When:
		ret = n.Short
	case Instead:
		ret = event.BeforePhase.PatternName(n.Short)
	case Report:
		ret = event.AfterPhase.PatternName(n.Short)
	default:
		ret = n.Prefix.String() + " " + n.Short
	}
	return
}

// pattern name as specified
// optional rule name as specified
// ex. Define rule:named:do: ["activating", "the standard activating action" ]
func ReadPhrase(patternSpec, ruleSpec string) (ret RuleName) {
	patternSpec = inflect.Normalize(patternSpec)
	name, suffix := findSuffix(patternSpec)
	// return name sans any prefix, and any prefix the name had.
	short, prefix := findPrefix(name)
	// fix: probably want some sort of "try prefix/suffix" that attempts to chop the parts
	// but restores them if it cant find them --
	// maybe see jess -- it does that sort of partial matching
	// and itd be neat to be able to use it here.
	var excludesPlayer bool
	const someone = "someone "
	if excludesPlayer = strings.HasPrefix(short, someone); excludesPlayer {
		short = short[len(someone):]
	}
	return RuleName{
		Short:  short,
		Label:  inflect.Normalize(ruleSpec),
		Prefix: prefix,
		Suffix: suffix,
	}
}

// match an author specified pattern reference
// to various naming conventions and pattern definitions
// to determine the intended pattern name, rank, and termination behavior.
// for example: "instead of x", "before x", "after x", "report x".
func (n RuleName) GetRuleInfo(run rt.Runtime) (ret RuleInfo, err error) {
	// fix: we can pass in the base type
	if ks, e := run.GetField(meta.KindAncestry, n.Short); e != nil {
		err = e
	} else {
		//
		switch pattern.Categorize(ks.Strings()) {
		default:
			err = errutil.Fmt("can't have a %q event", n.Short)

		// for regular patterns, supports sorting rules before/after
		case pattern.Calls:
			switch n.Prefix {
			case Instead, Report:
				err = errutil.Fmt("%q isn't an action and doesn't support %q", n.Short, n.Prefix)
			default:
				// ex. "before normal pattern", return "normal pattern"
				ret = RuleInfo{Name: n.Short, Rank: n.Prefix.rank()}
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
