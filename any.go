package anypb

import "encoding/json"

func (m *Any) IsEmpty() bool {
	if m != nil {
		switch m.GetType() {
		case Any_StringType:
			return (m.StringValue == nil || *m.StringValue == "")
		case Any_Uint32Type:
			return (m.Uint32Value == nil || *m.Uint32Value == 0)
		case Any_UintType, Any_Uint64Type:
			return (m.Uint64Value == nil || *m.Uint64Value == 0)
		case Any_Int32Type:
			return (m.Int32Value == nil || *m.Int32Value == 0)
		case Any_IntType, Any_Int64Type:
			return (m.Int64Value == nil || *m.Int64Value == 0)
		case Any_Float32Type:
			return (m.Float32Value == nil || *m.Float32Value == 0)
		case Any_Float64Type:
			return (m.Float64Value == nil || *m.Float64Value == 0)
		case Any_BoolType:
			return m.BoolValue == nil
		case Any_ByteType:
			return (m.ByteValue == nil || len(m.ByteValue) == 0)
		}
	}
	return true
}

func (m *Any) UnmarshalJSON(data []byte) error {
	var packet interface{}
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return err
	}
	*m = Any{}
	setAnyValue(packet, m)
	return nil
}
func (m *Any) MarshalJSON() ([]byte, error) {
	data := getAnyValue(m)
	return json.Marshal(data)
}

func UnmarshalAny(m *Any, v interface{}) error {
	var d decodeState
	data := getAnyValue(m)
	d.init(data)
	return d.unmarshal(v)
}
func MarshalAny(value interface{}) (m *Any, err error) {
	m = &Any{}
	err = setAnyValue(value, m)
	return
}

func (m *AnyMap) IsEmpty() bool {
	if m != nil {
		switch m.GetAnyType() {
		case AnyMap_AnyValueType:
			return m.GetAnyValue().IsEmpty()
		case AnyMap_AnyArrayType:
			return (m.AnyArray == nil || len(m.AnyArray) == 0)
		case AnyMap_AnyStringMapType, AnyMap_AnyFloat32MapType, AnyMap_AnyFloat64MapType:
			return (m.AnyStringMap == nil || len(m.AnyStringMap) == 0)
		case AnyMap_AnyUint32MapType:
			return (m.AnyUint32Map == nil || len(m.AnyUint32Map) == 0)
		case AnyMap_AnyUint64MapType, AnyMap_AnyUintMapType:
			return (m.AnyUint64Map == nil || len(m.AnyUint64Map) == 0)
		case AnyMap_AnyInt32MapType:
			return (m.AnyInt32Map == nil || len(m.AnyInt32Map) == 0)
		case AnyMap_AnyInt64MapType, AnyMap_AnyIntMapType:
			return (m.AnyInt64Map == nil || len(m.AnyInt64Map) == 0)
		case AnyMap_AnyBoolMapType:
			return (m.AnyBoolMap == nil || len(m.AnyBoolMap) == 0)
		}
	}
	return true
}

func (m *AnyMap) UnmarshalJSON(data []byte) error {
	var packet interface{}
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return err
	}
	*m = AnyMap{}
	setAnyMapValue(packet, m)
	return nil
}
func (m *AnyMap) MarshalJSON() ([]byte, error) {
	data := getAnyMapValue(m)
	return json.Marshal(data)
}

func UnmarshalAnyMap(m *AnyMap, v interface{}) error {
	var d decodeState
	data := getAnyMapValue(m)
	d.init(data)
	return d.unmarshal(v)
}
func MarshalAnyMap(value interface{}) (m *AnyMap, err error) {
	m = &AnyMap{}
	err = setAnyMapValue(value, m)
	return
}
