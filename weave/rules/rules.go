package rules

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

type RuleInfo struct {
	Name           string // name of the pattern / kind
	Rank           int    // smaller ranked rules run first
	Cancels        bool   // did the author intend the rule to terminate
	ExcludesPlayer bool
}

// match an author specified pattern reference
// to various naming conventions and pattern definitions
// to determine the intended pattern name, rank, and termination behavior.
// for example: "instead of x", "before x", "after x", "report x".
func ReadName(w g.Kinds, name string) (ret RuleInfo, err error) {
	short := name
	prefixIndex := numPrefixes // preliminary
	for i := 0; i < numPrefixes; i++ {
		p := eventPrefix(i).String()
		if strings.HasPrefix(name, p+" ") {
			short = name[len(p)+1:]
			prefixIndex = i
			break
		}
	}
	// fix: probably want some sort of "try prefix/suffix" that attempts to chop the parts
	// but restores them if it cant find them --
	// maybe see grok -- it does that sort of partial matching
	// and itd be neat to be able to use it here.
	var excludesPlayer bool
	const someone = "someone "
	if excludesPlayer = strings.HasPrefix(short, someone); excludesPlayer {
		next := short[len(someone):]
		if name == short {
			name = next
		}
		short = next
	}

	prefix := eventPrefix(prefixIndex)
	if k, e := w.GetKindByName(short); e != nil {
		err = e
	} else if events := kindsOf.Event.String(); k.Implements(events) {
		// block "before before traveling"
		// maybe eventually "first before traveling" or something like that.
		err = errutil.Fmt("can't have %q %s", name, events)
	} else if actions := kindsOf.Action.String(); !k.Implements(actions) {
		switch prefix {
		case instead, report:
			err = errutil.Fmt("%q isn't a kind of %s and doesn't support %q", short, actions, prefix)
		default:
			// ex. "before normal pattern", return "normal pattern"
			ret = RuleInfo{Name: short, Rank: prefix.rank()}
		}
	} else {
		var pattern string
		var cancels bool
		switch prefix {
		case instead:
			// ex. "instead of some action, return "before some action"
			pattern = event.BeforePhase.PatternName(short)
			cancels = true
		case report:
			// ex. "report some action, return "after some action"
			pattern = event.AfterPhase.PatternName(short)
			cancels = true
		case before:
			cancels = true
			fallthrough
		default:
			// ex. "before some pattern, return "before some pattern"
			pattern = name
		}
		ret = RuleInfo{
			Name:           pattern,
			Rank:           prefix.rank(),
			Cancels:        cancels,
			ExcludesPlayer: excludesPlayer,
		}
	}
	return
}

// Check to see if there are counters that might need updating on the regular.
// tdb: could this be processed at load time (storyImport)
func DoesUpdate(exe []rt.Execute) (okay bool) {
	// wrap the exes up into a single block for easier searching
	return searchCounters(&core.ChooseNothingElse{
		Exe: exe,
	})
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
	if op, e := searchForFlow(&core.ChooseNothingElse{
		Exe: exe,
	}, render.RenderResponse_Type); e != nil && e != jsn.Missing {
		panic(e)
	} else if response, ok := op.(*render.RenderResponse); ok && response.Text != nil {
		ret = response.Name
	}
	return
}
