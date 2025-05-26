package cueconvert

import "reflect"

type OutputValuer interface {
	OutputValues() map[string]any
}

type OutputFields map[string][]int

func (fields OutputFields) OutputValues(rv reflect.Value) map[string]any {
	if valuer, ok := rv.Interface().(OutputValuer); ok {
		return valuer.OutputValues()
	}

	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	values := map[string]any{}

	for name, loc := range fields {
		f := getField(rv, loc)

		if name == "" {
			if f.Kind() == reflect.Map {
				for _, k := range f.MapKeys() {
					key := k.String()
					if key == "$$task" {
						continue
					}
					values[key] = f.MapIndex(k).Interface()
				}
			}
			continue
		}

		// nil value never as output value
		if f.Kind() == reflect.Ptr {
			if !f.IsNil() {
				values[name] = f.Interface()
			}
			continue
		}

		values[name] = f.Interface()
	}

	return values
}

func getField(rv reflect.Value, loc []int) reflect.Value {
	switch len(loc) {
	case 0:
		return rv
	case 1:
		return rv.Field(loc[0])
	default:
		return getField(rv.Field(loc[0]), loc[1:])
	}
}
