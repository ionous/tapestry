package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

func getNewFloats(run rt.Runtime, assign core.Assignment) (ret []float64, err error) {
	if assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Number:
				if one, e := v.GetNumber(); e != nil {
					err = e
				} else {
					ret = []float64{one}
				}
			case affine.NumList:
				if many, e := v.GetNumList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.New("cant add %s to a num list", a)
			}
		}
	}
	return
}

func getNewStrings(run rt.Runtime, assign core.Assignment) (ret []string, err error) {
	if assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Text:
				if one, e := v.GetText(); e != nil {
					err = e
				} else {
					ret = []string{one}
				}
			case affine.TextList:
				if many, e := v.GetTextList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.New("cant add %s to a text list", a)
			}
		}
	}
	return
}
