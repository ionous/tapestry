package express

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type Render struct {
	Expression rt.TextEval `if:"internal"`
}

func (*Render) Compose() composer.Spec {
	return composer.Spec{
		Name:  "render_template",
		Spec:  "the template {lines%template:lines|quote}",
		Group: "format",
		Desc:  "Render Template: Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates",
	}
}

// RunTest returns an error on failure.
func (op *Render) GetText(run rt.Runtime) (ret string, err error) {
	return rt.GetText(run, op.Expression)
}
