package anypb

import (
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func getAnyValue(m *Any) interface{} {
	switch m.GetType() {
	case Any_StringType:
		return m.GetStringValue()
	case Any_UintType:
		return uint(m.GetUint64Value())
	case Any_Uint32Type:
		return m.GetUint32Value()
	case Any_Uint64Type:
		return m.GetUint64Value()
	case Any_IntType:
		return int(m.GetInt64Value())
	case Any_Int32Type:
		return m.GetInt32Value()
	case Any_Int64Type:
		return m.GetInt64Value()
	case Any_Float32Type:
		return m.GetFloat32Value()
	case Any_Float64Type:
		return m.GetFloat64Value()
	case Any_BoolType:
		return m.GetBoolValue()
	case Any_ByteType:
		return m.GetByteValue()
	}
	return nil
}

func getAnyMapValue(m *AnyMap) interface{} {
	switch m.GetAnyType() {
	case AnyMap_AnyValueType:
		return getAnyValue(m.GetAnyValue())
	case AnyMap_AnyArrayType:
		data := m.GetAnyArray()
		if data != nil {
			res := make([]interface{}, len(data))
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyStringMapType:
		data := m.GetAnyStringMap()
		if data != nil {
			res := make(map[string]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyUintMapType:
		data := m.GetAnyUint64Map()
		if data != nil {
			res := make(map[uint]interface{})
			for k, v := range data {
				res[uint(k)] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyUint32MapType:
		data := m.GetAnyUint32Map()
		if data != nil {
			res := make(map[uint32]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyUint64MapType:
		data := m.GetAnyUint64Map()
		if data != nil {
			res := make(map[uint64]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyIntMapType:
		data := m.GetAnyInt64Map()
		if data != nil {
			res := make(map[int]interface{})
			for k, v := range data {
				res[int(k)] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyInt32MapType:
		data := m.GetAnyInt32Map()
		if data != nil {
			res := make(map[int32]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyInt64MapType:
		data := m.GetAnyInt64Map()
		if data != nil {
			res := make(map[int64]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	case AnyMap_AnyFloat32MapType:
		data := m.GetAnyStringMap()
		if data != nil {
			res := make(map[float32]interface{})
			for k, v := range data {
				float, err := strconv.ParseFloat(k, 32)
				if err != nil {
					res[float32(float)] = getAnyMapValue(v)
				}
			}
			return res
		}
	case AnyMap_AnyFloat64MapType:
		data := m.GetAnyStringMap()
		if data != nil {
			res := make(map[float64]interface{})
			for k, v := range data {
				float, err := strconv.ParseFloat(k, 64)
				if err != nil {
					res[float] = getAnyMapValue(v)
				}
			}
			return res
		}
	case AnyMap_AnyBoolMapType:
		data := m.GetAnyBoolMap()
		if data != nil {
			res := make(map[bool]interface{})
			for k, v := range data {
				res[k] = getAnyMapValue(v)
			}
			return res
		}
	}
	return nil
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case runtime.Error:
				panic(r)
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d.value(rv)
	return d.savedError
}

type decodeState struct {
	data       interface{}
	savedError error
}

func (d *decodeState) init(data interface{}) *decodeState {
	d.data = data
	d.savedError = nil
	return d
}

func (d *decodeState) error(err error) {
	panic(err)
}

func (d *decodeState) saveError(err error) {
	if d.savedError == nil {
		d.savedError = err
	}
}

func (d *decodeState) value(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	d.next(v, reflect.ValueOf(d.data))
}

func (d *decodeState) next(v reflect.Value, sv reflect.Value) {
	pv := d.indirect(v, false)
	if pv.Kind() == reflect.Interface {
		pv.Set(sv)
		return
	}
	if sv.Kind() == reflect.Interface {
		sv = reflect.ValueOf(sv.Interface())
	}
	nv := d.indirect(sv, false)
	pt := pv.Type()
	if pv.Kind() == reflect.Struct {
		if nv.Kind() != reflect.Map {
			d.saveError(&UnmarshalTypeError{"struct", nv.Type()})
			return
		}

		nt := nv.Type()
		if nt.Key().Kind() != reflect.String {
			d.saveError(&UnmarshalTypeError{"struct", nv.Type()})
			return
		}

		fields := make(map[string]string)
		l := pt.NumField()
		for i := 0; i < l; i++ {
			pf := pt.Field(i)
			tag := pf.Tag.Get("anypb")
			if tag != "" {
				ptp := strings.SplitN(tag, ",", 2)
				fields[ptp[0]] = pf.Name
			} else {
				fields[pf.Name] = pf.Name
			}
		}

		keys := nv.MapKeys()
		for _, k := range keys {
			n, ok := fields[k.String()]
			if !ok {
				d.saveError(&UnmarshalTypeError{"struct", nv.Type()})
				return
			}
			mv := pv.FieldByName(n)
			d.next(mv, nv.MapIndex(k))
		}
		return
	}
	if pv.Kind() != nv.Kind() {
		d.saveError(&UnmarshalTypeError{pv.Type().String(), nv.Type()})
		return
	}
	switch pv.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if pv.IsNil() && pv.Kind() == reflect.Slice {
			pv.Set(reflect.MakeSlice(pt, 0, 0))
		}
		l := nv.Len()
		for i := 0; i < l; i++ {
			if pv.Kind() == reflect.Slice {
				if i >= pv.Cap() {
					newcap := pv.Cap() + pv.Cap()/2
					if newcap < 4 {
						newcap = 4
					}
					newv := reflect.MakeSlice(pv.Type(), pv.Len(), newcap)
					reflect.Copy(newv, pv)
					pv.Set(newv)
				}
				if i >= pv.Len() {
					pv.SetLen(i + 1)
				}
			}
			if i < pv.Len() {
				d.next(pv.Index(i), nv.Index(i))
			}
		}
		if l < pv.Len() {
			z := reflect.Zero(pv.Type().Elem())
			if pv.Kind() == reflect.Array {
				for i := l; i < pv.Len(); i++ {
					pv.Index(i).Set(z)
				}
			} else {
				pv.SetLen(l)
			}
		}
	case reflect.Map:
		if pv.IsNil() {
			pv.Set(reflect.MakeMap(pt))
		}
		nt := nv.Type()
		if pt.Key().Kind() != nt.Key().Kind() {
			d.saveError(&UnmarshalTypeError{"map", nv.Type()})
			return
		}
		keys := nv.MapKeys()
		for _, k := range keys {
			mv := pv.MapIndex(k)
			if !mv.IsValid() {
				mv = reflect.New(pt.Elem())
				pv.SetMapIndex(k, mv)
			}
			d.next(mv, nv.MapIndex(k))
		}
	default:
		pv.Set(nv)
	}
}

func (d *decodeState) indirect(v reflect.Value, decodingNull bool) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "anypb: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "anypb: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "anypb: Unmarshal(nil " + e.Type.String() + ")"
}

type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
