package raw

import (
	"cmp"
	"fmt"
	"slices"

	"git.sr.ht/~ionous/tapestry/rt"
)

func FindKind(ks []rt.Kind, exactKind string) (ret *rt.Kind, err error) {
	if i, ok := slices.BinarySearchFunc(ks, exactKind, func(k rt.Kind, _ string) int {
		return cmp.Compare(k.Name(), exactKind)
	}); !ok {
		err = fmt.Errorf("couldnt find kind %q", exactKind)
	} else {
		ret = &ks[i]
	}
	return
}
