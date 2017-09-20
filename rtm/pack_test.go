package rtm_test

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	. "github.com/ionous/iffy/tests"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// TestRegistration tests the object builder.
func TestRegistration(t *testing.T) {
	assert := testify.New(t)
	first := &BaseClass{Name: "first", State: Yes, Labeled: true}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe}}
	run := newRuntime(first, second)
	if n, ok := run.GetObject("first"); assert.True(ok) {
		assert.Equal(ident.IdOf("$first"), n.Id())
	}
	if d, ok := run.GetObject("second"); assert.True(ok) {
		assert.Equal(ident.IdOf("$second"), d.Id())
	}
}

// TestStateAccess
func TestStateAccess(t *testing.T) {
	assert := testify.New(t)
	first := &BaseClass{Name: "first", State: Yes, Labeled: true}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe}}
	//
	run := newRuntime(first, second)
	test := func(name, prop string, value bool) {
		if obj, ok := run.GetObject(name); assert.True(ok) {
			var res bool
			if e := obj.GetValue(prop, &res); assert.NoError(e) {
				if !assert.Equal(value, res) {
					t.Log("mismatched", obj, prop)
				}
			}
		}
	}

	test("first", "yes", true)
	test("first", "no", false)
	test("first", "maybe", false)
	test("first", "labeled", true)
	//
	test("second", "yes", false)
	test("second", "no", false)
	test("second", "maybe", true)
	test("second", "labeled", false)
}

func TestStateSet(t *testing.T) {
	assert := testify.New(t)

	first := &BaseClass{Name: "first", State: Yes, Labeled: true, Object: ident.IdOf("second")}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe, Object: ident.IdOf("first")}}

	run := newRuntime(first, second)
	unpackValue := func(obj rt.Object, name string, pv interface{}) {
		if e := obj.GetValue(name, pv); e != nil {
			panic(e)
		}
	}
	packValue := func(obj rt.Object, name string, v interface{}) {
		if e := obj.SetValue(name, v); e != nil {
			panic(e)
		}
	}

	if n, ok := run.GetObject("first"); assert.True(ok) {
		var res bool
		// start with yes, it should be true
		unpackValue(n, "yes", &res)
		if assert.True(res) {
			// try to change the value to maybe
			packValue(n, "maybe", true)
			// yes should now be false.
			unpackValue(n, "yes", &res)
			if assert.False(res) {
				// and maybe should now be true
				unpackValue(n, "maybe", &res)
				assert.True(res)
				// try to change states in an illegal way:
				e := n.SetValue("maybe", false)
				assert.Error(e)

				// add verify it didnt change:
				unpackValue(n, "maybe", &res)
				assert.True(res)
			}
		}
		//
		t.Run("string", func(t *testing.T) {
			unpackValue(n, "yes", &res)
			if res {
				t.Fatal("yes should be false")
			} else {
				packValue(n, "state", "yes")

				unpackValue(n, "yes", &res)
				if !res {
					t.Fatal("yes should be true")
				}
			}
		})
	}
	// check, change, and check the labeled bool.
	toggle := func(name, prop string, goal bool) {
		if n, ok := run.GetObject(name); assert.True(ok) {
			var res bool
			unpackValue(n, prop, &res)
			if assert.NotEqual(goal, res, "initial value") {
				packValue(n, prop, goal)
				unpackValue(n, prop, &res)
				assert.Equal(goal, res)
			}
		}
	}
	toggle("second", "labeled", true)
	toggle("second", "labeled", false)
}

func newRuntime(ptrs ...interface{}) rt.Runtime {
	var objects obj.Registry
	objects.RegisterValues(ptrs)
	run, e := rtm.New(nil).Objects(objects).Rtm()
	if e != nil {
		panic(e)
	}
	return run
}

// TestPropertyAccess to ensure normal properties are accessible
func TestPropertyAccess(t *testing.T) {
	first := &BaseClass{Name: "first", State: Yes, Labeled: true, Object: ident.IdOf("second")}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe, Object: ident.IdOf("first")}}
	run := newRuntime(first, second)

	// we create some slots for values to be unpacked into
	var expected = []struct {
		name string
		pv   interface{}
	}{
		{"Name", new(string)},
		{"Num", new(float64)},
		{"Text", new(string)},
		{"Object", new(rt.Object)},
		{"Nums", new([]float64)},
		{"Texts", new([]string)},
		{"Objects", new([]rt.Object)},
	}
	test := func(n rt.Object) (err error) {
		for _, v := range expected {
			if e := n.GetValue(v.name, v.pv); e != nil {
				err = e
				break
			}
		}
		return
	}
	//
	assert := testify.New(t)
	if n, ok := run.GetObject("first"); assert.True(ok) {
		if d, ok := run.GetObject("second"); assert.True(ok) {
			// from n get d:
			if e := test(n); assert.NoError(e) {
				other := (*expected[3].pv.(*rt.Object))
				if assert.Equal(d, other) {
					// from d get n:
					if e := test(d); assert.NoError(e) {
						other := (*expected[3].pv.(*rt.Object))
						assert.Equal(n, other)
					}
				}
			}
		}
	}
}
