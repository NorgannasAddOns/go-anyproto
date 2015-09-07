package anypb

import (
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

func setData(anyVal reflect.Value, anyType Any_Type, m *Any, val reflect.Value) (err error) {
	m.Type = &anyType
	anyVal.Elem().Set(val)
	return
}
func setAnyValue(value interface{}, m *Any) (err error) {
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
	if m != nil {
		if m.Type != nil {
			m.Reset()
		}
		var val reflect.Value
		if v, ok := value.(reflect.Value); ok {
			val = v
		} else if v, ok := value.(*Any); ok {
			if v != nil {
				*m = *v
			} else {
				anyType := Any_NilType
				m.Type = &anyType
			}
			return
		} else {
			val = reflect.ValueOf(value)
		}
		if !val.IsValid() {
			anyType := Any_NilType
			m.Type = &anyType
			return
		}
		kind := val.Kind()
		if (kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface || kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice) && val.IsNil() {
			anyType := Any_NilType
			m.Type = &anyType
			return
		}
		if kind == reflect.Ptr {
			kind = val.Elem().Kind()
		}
		switch kind {
		case reflect.Interface:
			err = setAnyValue(val.Interface(), m)
		case reflect.String:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.String(val.String()))
			}
			err = setData(reflect.ValueOf(&m.StringValue), Any_StringType, m, val)
		case reflect.Uint:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Uint64(val.Uint()))
			}
			err = setData(reflect.ValueOf(&m.Uint64Value), Any_UintType, m, val)
		case reflect.Uint32:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Uint32(uint32(val.Uint())))
			}
			err = setData(reflect.ValueOf(&m.Uint32Value), Any_Uint32Type, m, val)
		case reflect.Uint64:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Uint64(val.Uint()))
			}
			err = setData(reflect.ValueOf(&m.Uint64Value), Any_Uint64Type, m, val)
		case reflect.Int:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Int64(val.Int()))
			}
			err = setData(reflect.ValueOf(&m.Int64Value), Any_IntType, m, val)
		case reflect.Int32:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Int32(int32(val.Int())))
			}
			err = setData(reflect.ValueOf(&m.Int32Value), Any_Int32Type, m, val)
		case reflect.Int64:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Int64(val.Int()))
			}
			err = setData(reflect.ValueOf(&m.Int64Value), Any_Int64Type, m, val)
		case reflect.Float32:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Float32(float32(val.Float())))
			}
			err = setData(reflect.ValueOf(&m.Float32Value), Any_Float32Type, m, val)
		case reflect.Float64:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Float64(val.Float()))
			}
			err = setData(reflect.ValueOf(&m.Float64Value), Any_Float64Type, m, val)
		case reflect.Bool:
			if kind != reflect.Ptr {
				val = reflect.ValueOf(proto.Bool(val.Bool()))
			}
			err = setData(reflect.ValueOf(&m.BoolValue), Any_BoolType, m, val)
		case reflect.Slice:
			if val.Type() == typeOfBytes {
				err = setData(reflect.ValueOf(&m.ByteValue), Any_ByteType, m, val)
				break
			}
			err = errors.New("Error: Unsupported value type")
		default:
			err = errors.New("Error: Unsupported value type")
		}
	}
	return
}

func setMapData(kind reflect.Kind, data reflect.Value, anyMap reflect.Value, anyType AnyMap_Type, m *AnyMap, val reflect.Value, keys []reflect.Value, keyFn func(key reflect.Value) reflect.Value) (err error) {
	mapData := data.Elem()
	for _, key := range keys {
		if key.Kind() == reflect.Ptr {
			key = key.Elem()
		}
		if key.Kind() != kind {
			err = errors.New("Error: Map key types don't match")
			return
		}
		any := &AnyMap{}
		err = setAnyMapValue(val.MapIndex(key), any)
		if err != nil {
			return
		}
		if keyFn != nil {
			key = keyFn(key)
		}
		mapData.SetMapIndex(key, reflect.ValueOf(any))
	}
	m.AnyType = &anyType
	anyMap.Elem().Set(mapData)
	return
}
func setAnyMapValue(value interface{}, m *AnyMap) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	if m != nil {
		if m.AnyType != nil {
			m.Reset()
		}
		var val reflect.Value
		if v, ok := value.(reflect.Value); ok {
			val = v
		} else if v, ok := value.(*AnyMap); ok {
			if v != nil {
				*m = *v
			} else {
				anyType := AnyMap_NilType
				m.AnyType = &anyType
			}
			return
		} else if v, ok := value.(*Any); ok {
			if v != nil {
				anyType := AnyMap_AnyValueType
				m.AnyType = &anyType
				m.AnyValue = v
			} else {
				anyType := AnyMap_NilType
				m.AnyType = &anyType
			}
			return
		} else {
			val = reflect.ValueOf(value)
		}
		if !val.IsValid() {
			anyType := AnyMap_NilType
			m.AnyType = &anyType
			return
		}
		kind := val.Kind()
		if (kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface || kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice) && val.IsNil() {
			anyType := AnyMap_NilType
			m.AnyType = &anyType
			return
		}
		if kind == reflect.Ptr {
			kind = val.Elem().Kind()
		}
		switch kind {
		case reflect.Interface:
			err = setAnyMapValue(val.Interface(), m)
		case reflect.Slice:
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.Type() == typeOfBytes {
				anyType := AnyMap_AnyValueType
				m.AnyType = &anyType
				m.AnyValue = &Any{}
				err = setAnyValue(val, m.AnyValue)
				break
			}
			if val.Len() > 0 {
				anyType := AnyMap_AnyArrayType
				m.AnyType = &anyType
				m.AnyArray = make([]*AnyMap, val.Len())
				for i := 0; i < val.Len(); i++ {
					any := &AnyMap{}
					setAnyMapValue(val.Index(i), any)
					m.AnyArray[i] = any
				}
			}
		case reflect.Struct:
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			data := make(map[string]*AnyMap)
			t := val.Type()
			l := t.NumField()
			for i := 0; i < l; i++ {
				f := t.Field(i)
				v := val.FieldByName(f.Name)
				tag := f.Tag.Get("anypb")
				if tag != "" {
					tp := strings.SplitN(tag, ",", 2)
					if tp[0] == "-" {
						continue
					}
					any := &AnyMap{}
					setAnyMapValue(v, any)
					if len(tp) > 1 && strings.Contains(tp[1], "omitempty") && any.IsEmpty() {
						continue
					}
					data[tp[0]] = any
				} else {
					any := &AnyMap{}
					setAnyMapValue(v, any)
					if any.IsEmpty() {
						continue
					}
					data[f.Name] = any
				}
			}
			anyType := AnyMap_AnyStringMapType
			m.AnyType = &anyType
			m.AnyStringMap = data
		case reflect.Map:
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			keys := val.MapKeys()
			if len(keys) > 0 {
				kind := keys[0].Kind()
				if kind == reflect.Ptr {
					kind = keys[0].Elem().Kind()
				}
				switch kind {
				case reflect.String:
					data := make(map[string]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyStringMap), AnyMap_AnyStringMapType, m, val, keys, nil)
				case reflect.Uint:
					data := make(map[uint64]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyUint64Map), AnyMap_AnyUintMapType, m, val, keys, nil)
				case reflect.Uint32:
					data := make(map[uint32]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyUint32Map), AnyMap_AnyUint32MapType, m, val, keys, nil)
				case reflect.Uint64:
					data := make(map[uint64]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyUint64Map), AnyMap_AnyUint64MapType, m, val, keys, nil)
				case reflect.Int:
					data := make(map[int64]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyInt64Map), AnyMap_AnyIntMapType, m, val, keys, nil)
				case reflect.Int32:
					data := make(map[int32]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyInt32Map), AnyMap_AnyInt32MapType, m, val, keys, nil)
				case reflect.Int64:
					data := make(map[int64]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyInt64Map), AnyMap_AnyInt64MapType, m, val, keys, nil)
				case reflect.Float32:
					data := make(map[string]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyStringMap), AnyMap_AnyFloat32MapType, m, val, keys, func(key reflect.Value) reflect.Value {
						return reflect.ValueOf(strconv.FormatFloat(key.Float(), 'f', -1, 32))
					})
				case reflect.Float64:
					data := make(map[string]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyStringMap), AnyMap_AnyFloat64MapType, m, val, keys, func(key reflect.Value) reflect.Value {
						return reflect.ValueOf(strconv.FormatFloat(key.Float(), 'f', -1, 64))
					})
				case reflect.Bool:
					data := make(map[bool]*AnyMap)
					err = setMapData(kind, reflect.ValueOf(&data), reflect.ValueOf(&m.AnyBoolMap), AnyMap_AnyBoolMapType, m, val, keys, nil)
				default:
					err = errors.New("Error: Unsupported map key type")
				}
			}
		default:
			anyType := AnyMap_AnyValueType
			m.AnyType = &anyType
			m.AnyValue = &Any{}
			err = setAnyValue(val, m.AnyValue)
		}
	}
	return
}
