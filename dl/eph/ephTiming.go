package eph

import "git.sr.ht/~ionous/tapestry/imp/assert"

func toTiming(t EphTiming, always EphAlways) (ret assert.EventTiming) {
	switch t.Str {
	case EphTiming_Before:
		ret = assert.Before
	case EphTiming_During:
		ret = assert.During
	case EphTiming_After:
		ret = assert.After
	case EphTiming_Later:
		ret = assert.Later
	}
	if always.Str == EphAlways_Always {
		ret |= assert.RunAlways
	}
	return
}

func fromTiming(t assert.EventTiming) (ret EphTiming, always EphAlways) {
	if t&assert.RunAlways != 0 {
		always.Str = EphAlways_Always
		t ^= assert.RunAlways
	}
	switch t {
	case assert.Before:
		ret.Str = EphTiming_Before
	case assert.During:
		ret.Str = EphTiming_During
	case assert.After:
		ret.Str = EphTiming_After
	case assert.Later:
		ret.Str = EphTiming_Later
	}
	return
}
