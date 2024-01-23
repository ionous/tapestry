package rules

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	inflect "git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"github.com/ionous/errutil"
)

type RuleInfo struct {
	Name           string // name of the pattern / kind
	Label          string
	Rank           int // smaller ranked rules run first
	Stop           bool
	Jump           rt.Jump
	ExcludesPlayer bool // true if the rule should apply to all actors
}

type RuleName struct {
	Short, Long, Label string
	Prefix             Prefix
	Suffix             Suffix
	ExcludesPlayer     bool // true if the rule should apply to all actors
}

func (n RuleName) IsDomainEvent() bool {
	return n.Suffix == Begins || n.Suffix == Ends
}

// instead and report are grouped with before and after respectively
func (n RuleName) EventName() (ret string) {
	if n.IsDomainEvent() {
		var s strings.Builder
		s.WriteString(n.Short)
		s.WriteRune(' ')
		s.WriteString(n.Suffix.String())
		ret = s.String()
	} else {
		switch n.Prefix {
		case When:
			ret = n.Short
		case Instead:
			ret = event.BeforePhase.PatternName(n.Short)
		case Report:
			ret = event.AfterPhase.PatternName(n.Short)
		default:
			ret = n.Long
		}
	}
	return
}

func ReadPhrase(phrase, label string) (ret RuleName) {
	phrase, label = inflect.Normalize(phrase), inflect.Normalize(label)
	name, suffix := findSuffix(phrase)
	// return name sans any prefix, and any prefix the name had.
	short, prefix := findPrefix(name)
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
	return RuleName{
		Short:  short,
		Long:   name,
		Label:  label,
		Prefix: prefix,
		Suffix: suffix,
	}
}

func (n RuleName) GetRuleInfo(ks g.Kinds) (ret RuleInfo, err error) {
	if n.IsDomainEvent() {
		ret, err = n.ruleForDomain()
	} else {
		ret, err = n.ruleForPattern(ks)
	}
	return
}

func (n RuleName) ruleForDomain() (ret RuleInfo, err error) {
	ret = RuleInfo{
		Name: n.EventName(),
		Jump: rt.JumpLater,
	}
	return
}

// match an author specified pattern reference
// to various naming conventions and pattern definitions
// to determine the intended pattern name, rank, and termination behavior.
// for example: "instead of x", "before x", "after x", "report x".
func (n RuleName) ruleForPattern(ks g.Kinds) (ret RuleInfo, err error) {
	if k, e := ks.GetKindByName(n.Short); e != nil {
		err = e // ^ the base pattern
	} else {
		switch pattern.Categorize(k) {
		default:
			err = errutil.Fmt("can't have a %q event", n.Short)

		case pattern.Calls:
			switch n.Prefix {
			case Instead, Report:
				err = errutil.Fmt("%q isn't a kind of %s and doesn't support %q", n.Short, kindsOf.Action, n.Prefix)
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
	if op, e := searchForFlow(&m, render.RenderResponse_Type); e != nil && e != jsn.Missing {
		panic(e)
	} else if response, ok := op.(*render.RenderResponse); ok && response.Text != nil {
		ret = response.Name
	}
	return
}
