package phpserialize

import (
	"fmt"
	"reflect"
)

var valueEncoders []encoderFunc

//nolint:gochecknoinits
func init() {
	valueEncoders = []encoderFunc{
		reflect.Bool:          encodeBoolValue,
		reflect.Int:           encodeIntValue,
		reflect.Int8:          encodeIntValue,
		reflect.Int16:         encodeIntValue,
		reflect.Int32:         encodeIntValue,
		reflect.Int64:         encodeIntValue,
		reflect.Uint:          encodeUintValue,
		reflect.Uint8:         encodeUintValue,
		reflect.Uint16:        encodeUintValue,
		reflect.Uint32:        encodeUintValue,
		reflect.Uint64:        encodeUintValue,
		reflect.Float32:       encodeFloat32Value,
		reflect.Float64:       encodeFloat64Value,
		reflect.Complex64:     encodeUnsupportedValue,
		reflect.Complex128:    encodeUnsupportedValue,
		reflect.Array:         encodeUnsupportedValue, //encodeArrayValue,
		reflect.Chan:          encodeUnsupportedValue,
		reflect.Func:          encodeUnsupportedValue,
		reflect.Interface:     encodeUnsupportedValue, //encodeInterfaceValue,
		reflect.Map:           encodeMapValue,
		reflect.Ptr:           encodeUnsupportedValue,
		reflect.Slice:         encodeSliceValue,
		reflect.String:        encodeStringValue,
		reflect.Struct:        encodeStructValue,
		reflect.UnsafePointer: encodeUnsupportedValue,
	}
}

func getEncoder(typ reflect.Type) encoderFunc {
	if v, ok := typeEncMap.Load(typ); ok {
		return v.(encoderFunc)
	}
	fn := _getEncoder(typ)
	if fn == nil {
		return encodeUnsupportedValue
	}
	typeEncMap.Store(typ, fn)
	return fn
}

func _getEncoder(typ reflect.Type) encoderFunc {
	kind := typ.Kind()

	if kind == reflect.Ptr {
		if _, ok := typeEncMap.Load(typ.Elem()); ok {
			return encodeUnsupportedValue //ptrEncoderFunc(typ)
		}
	}

	/*if typ.Implements(customEncoderType) {
		return encodeCustomValue
	}
	if typ.Implements(marshalerType) {
		return marshalValue
	}
	if typ.Implements(binaryMarshalerType) {
		return marshalBinaryValue
	}
	if typ.Implements(textMarshalerType) {
		return marshalTextValue
	}

	// Addressable struct field value.
	if kind != reflect.Ptr {
		ptr := reflect.PtrTo(typ)
		if ptr.Implements(customEncoderType) {
			return encodeCustomValuePtr
		}
		if ptr.Implements(marshalerType) {
			return marshalValuePtr
		}
		if ptr.Implements(binaryMarshalerType) {
			return marshalBinaryValueAddr
		}
		if ptr.Implements(textMarshalerType) {
			return marshalTextValueAddr
		}
	}*/

	/*if typ == errorType {
		return encodeErrorValue
	}*/

	/*switch kind {
	case reflect.Ptr:
		return ptrEncoderFunc(typ)
	case reflect.Slice:
		elem := typ.Elem()
		if elem.Kind() == reflect.Uint8 {
			return encodeByteSliceValue
		}
		if elem == stringType {
			return encodeStringSliceValue
		}
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return encodeByteArrayValue
		}
	case reflect.Map:
		if typ.Key() == stringType {
			switch typ.Elem() {
			case stringType:
				return encodeMapStringStringValue
			case interfaceType:
				return encodeMapStringInterfaceValue
			}
		}
	}
	*/

	return valueEncoders[kind]
}

func encodeUnsupportedValue(e *Encoder, v reflect.Value) error {
	return fmt.Errorf("phpserialize: Encode(unsupported %s)", v.Type())
}

func encodeBoolValue(e *Encoder, v reflect.Value) error {
	return e.EncodeBool(v.Bool())
}

func encodeIntValue(e *Encoder, v reflect.Value) error {
	return e.EncodeInt64(v.Int())
}

func encodeUintValue(e *Encoder, v reflect.Value) error {
	return e.EncodeUint64(v.Uint())
}

func encodeStringValue(e *Encoder, v reflect.Value) error {
	return e.EncodeString(v.String())
}

func encodeFloat32Value(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat64(v.Float())
}

func encodeFloat64Value(e *Encoder, v reflect.Value) error {
	return e.EncodeFloat64(v.Float())
}
