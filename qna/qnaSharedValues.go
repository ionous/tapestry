package qna

// // check the cache before asking the database for info
// // the cache stores "assignments" so we can then take a snapshot of the dynamic value at the current moment in time
// func (run *Runner) getDynamicValue(target, field string, queryFn queryFn) (g.Value, error) {
// 	key := makeKey(target, field)
// 	assign := run.values.getOrCache(key, queryFn)
// 	return assign.GetAssignedValue(n)
// }

// // check the cache before asking the database for info.
// func (n *sharedValues) getSharedValue(key keyType, sharedQuery sharedQuery) (ret g.Value, err error) {
// 	if x, ok := (*n)[key]; ok {
// 		ret, err = x.v, x.e
// 	} else {
// 		switch val, e := sharedQuery(); e {
// 		case nil: // success!
// 			ret = val
// 			(*n)[key] = sharedValue{v: ret}

// 		case sql.ErrNoRows: // no data.
// 			err = key.unknown()
// 			(*n)[key] = sharedValue{e: err}

// 		default: // some other error; dont cache.
// 			err = e
// 		}
// 	}
// 	return
// }

// func readValue(aff affine.Affinity, val interface{}) (ret g.Value, err error) {
// 	switch aff {
// 	case affine.Bool:
// 		switch v := val.(type) {
// 		case nil:
// 			ret = g.False // zero value for unhandled defaults in sqlite
// 		case bool:
// 			ret = g.BoolOf(v)
// 		case int64:
// 			// sqlite, boolean values can be represented as 1/0
// 			ret = g.BoolOf(v == 0)
// 		default:
// 			err = invalidValue(aff, val)
// 		}

// 	case affine.Number:
// 		switch v := val.(type) {
// 		case nil:
// 			ret = g.Zero // zero value for unhandled defaults in sqlite
// 		case int64:
// 			ret = g.IntOf(int(v))
// 		case float64:
// 			ret = g.FloatOf(v)
// 		default:
// 			err = invalidValue(aff, val)
// 		}

// 	case affine.Text:
// 		switch v := val.(type) {
// 		case nil:
// 			ret = g.Empty // zero value for unhandled defaults in sqlite
// 		case string:
// 			ret = g.StringOf(v)
// 		default:
// 			err = invalidValue(aff, val)
// 		}

// 	case affine.NumList:
// 		// FIX: fine for the second, but im assuming that actually we'll need to decode a comma-separated list of strings
// 		switch vs := i.(type) {
// 		case []float64:
// 			ret = g.FloatsOf(vs)
// 		default:
// 			err = invalidValue(aff, val)
// 		}

// 	case affine.TextList:
// 		switch vs := val.(type) {
// 		case []string:
// 			ret = g.StringsOf(vs)
// 		default:
// 			err = invalidValue(aff, val)
// 		}

// 	case affine.Record:
// 		// FIX: fine for the second, but im assuming that actually we'll need to decode a structured value
// 		// or unmarshal a make
// 		if v, ok := val.(*g.Record); ok {
// 			ret = g.RecordOf(v)
// 		} else {
// 			err = invalidValue(aff, val)
// 		}

// 	// we could either disallow direct record list storage, or:
// 	// store the requested kind for storage.
// 	// case affine.RecordList:
// 	// 	switch vs := val.(type) {
// 	// 	case []*g.Record:
// 	// 		ret = g.RecordsOf(vs)
// 	// 	default:
// 	// 	 err = invalidValue(aff, val)
// 	// 	}

// 	default:
// 		err = invalidValue("", val)
// 	}
// 	return
// }

// func invalidValue(aff affine.Affinity, val interface{}) error {
// 	return errutil.Fmt("make value: %T can't be converted to %s", val, aff.String())
// }
