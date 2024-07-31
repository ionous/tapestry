package format

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// alt: could collect into a private member of print menu choices
// alt: the parser could examine the text output and capture the choices.
// ( examining the ol / ul for numbers i suppose )
var CurrentMenu MenuData

func (op *PrintMenu) Execute(run rt.Runtime) (err error) {
	// alt: could use a single action and multiple scenes?
	// or: a single action, and some text passed as a filter.
	if action, e := safe.GetText(run, op.ActionName); e != nil {
		err = cmd.Error(op, e)
	} else if data, e := GetMenuData(run, action.String(), op.MenuOptions); e != nil {
		err = cmd.Error(op, e)
	} else {
		CurrentMenu = data
		if CurrentMenu.UseList {
			w := run.Writer()
			tag := "ul"
			if CurrentMenu.UseNumbers {
				tag = "ol"
			}
			fmt.Fprintf(w, "<%s>", tag)
			defer fmt.Fprintf(w, "</%s>", tag)
		}
		// generate the contents
		if e := safe.RunAll(run, op.Exe); e != nil {
			err = cmd.Error(op, e)
		} else {
			// alt: could be sent via Runtime.Call
			// or copied into a predefined meta.Field.
			// could implement a "thunk" for the runtime so it can be more like a vtable
			// and allow playtime to wrap it, check for an interface implementation here.
			if len(CurrentMenu.Choices) == 0 {
				log.Println("no menu options existed for", action.String())
			}
		}
	}
	return
}

func (op *PrintMenuChoice) Execute(run rt.Runtime) (err error) {
	if label, e := safe.GetText(run, op.Label); e != nil {
		err = cmd.Error(op, e)
	} else if content, e := safe.GetOptionalText(run, op.Content, ""); e != nil {
		err = cmd.Error(op, e)
	} else {
		w := run.Writer()
		if CurrentMenu.UseList {
			fmt.Fprintf(w, "<li>")
			defer fmt.Fprintf(w, "</li>")
		}
		fmt.Fprintf(w, "<a=%q>%s</a>", label.String(), content.String())
		CurrentMenu.Choices = append(CurrentMenu.Choices, label.String())
	}
	return
}

func GetMenuData(run rt.Runtime, name string, op *MenuOptions) (ret MenuData, err error) {
	ret = MenuData{
		Action:     name,
		UseList:    true,
		UseNumbers: true,
	}
	if op != nil {
		options := []struct {
			eval rt.BoolEval
			pval *bool
		}{
			{op.ShowList, &ret.UseList},
			{op.ShowNumbers, &ret.UseNumbers},
		}
		for _, opt := range options {
			if res, e := safe.GetOptionalBool(run, opt.eval, *opt.pval); e != nil {
				err = e
				break
			} else {
				(*opt.pval) = res.Bool()
			}
		}
	}
	return
}

type MenuData struct {
	Action string
	UseList,
	UseNumbers bool
	Choices []string
}

func (pt *MenuData) Match(w string) (ret string, okay bool) {
	if a, ok := pt.tryNumber(w); ok {
		ret, okay = a, true
	} else if a, ok := pt.tryChoice(w); ok {
		ret, okay = a, true
	}
	return
}

func (pt *MenuData) tryNumber(w string) (ret string, okay bool) {
	if pt.UseNumbers {
		if i, e := strconv.Atoi(w); e == nil {
			if i > 0 && i <= len(pt.Choices) {
				ret, okay = pt.Choices[i-1], true
			}
		}
	}
	return
}

func (pt *MenuData) tryChoice(w string) (ret string, okay bool) {
	// optional?
	if i := slices.Index(pt.Choices, w); i >= 0 {
		ret, okay = pt.Choices[i], true
	}
	return
}
