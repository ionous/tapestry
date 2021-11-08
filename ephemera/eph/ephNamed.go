package eph

import (
	"database/sql/driver"
	"encoding/json"
)

// opaque row id for name
// FIX: this was probably a mistake, and we should just store strings instead
// im not sure what value it really adds to pre-consolidate names
// especially since assembly has to unpack them all to do any comparisons.
type Named struct {
	id  int64
	str string
}

func MakeName(id int64, str string) Named {
	return Named{id, str}
}

func (ns *Named) IsValid() bool {
	return ns.id > 0
}

func (ns *Named) String() string {
	return ns.str
}

func (ns Named) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.str)
}

// Scan converts a database value into a Named entry. ( opposite of Value )
func (ns *Named) Scan(value interface{}) (err error) {
	if v, e := driver.DefaultParameterConverter.ConvertValue(value); e != nil {
		err = e
	} else {
		ns.id = v.(int64)
	}
	return
}

// Value converts a Named entry into a database value. ( opposite of Scan )
func (ns Named) Value() (driver.Value, error) {
	return ns.id, nil
}
