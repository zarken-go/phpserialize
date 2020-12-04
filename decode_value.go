package phpserialize

import (
	"fmt"
	"math/bits"
	"reflect"
)

func getDecoder(typ reflect.Type) decoderFunc {
	if v, ok := typeDecMap.Load(typ); ok {
		return v.(decoderFunc)
	}
	fn := _getDecoder(typ)
	typeDecMap.Store(typ, fn)
	return fn
}

func _getDecoder(typ reflect.Type) decoderFunc {
	kind := typ.Kind()

	if kind == reflect.Ptr {
		if _, ok := typeDecMap.Load(typ.Elem()); ok {
			return ptrDecoderFunc(typ)
		}
	}

	/*
		if typ.Implements(customDecoderType) {
			return decodeCustomValue
		}
		if typ.Implements(unmarshalerType) {
			return unmarshalValue
		}
		if typ.Implements(binaryUnmarshalerType) {
			return unmarshalBinaryValue
		}
		if typ.Implements(textUnmarshalerType) {
			return unmarshalTextValue
		}

		// Addressable struct field value.
		if kind != reflect.Ptr {
			ptr := reflect.PtrTo(typ)
			if ptr.Implements(customDecoderType) {
				return decodeCustomValueAddr
			}
			if ptr.Implements(unmarshalerType) {
				return unmarshalValueAddr
			}
			if ptr.Implements(binaryUnmarshalerType) {
				return unmarshalBinaryValueAddr
			}
			if ptr.Implements(textUnmarshalerType) {
				return unmarshalTextValueAddr
			}
		}

		switch kind {
		case reflect.Ptr:
			return ptrDecoderFunc(typ)
		case reflect.Slice:
			elem := typ.Elem()
			if elem.Kind() == reflect.Uint8 {
				return decodeBytesValue
			}
			if elem == stringType {
				return decodeStringSliceValue
			}
		case reflect.Array:
			if typ.Elem().Kind() == reflect.Uint8 {
				return decodeByteArrayValue
			}
		case reflect.Map:
			if typ.Key() == stringType {
				switch typ.Elem() {
				case stringType:
					return decodeMapStringStringValue
				case interfaceType:
					return decodeMapStringInterfaceValue
				}
			}
		}
	*/

	return valueDecoders[kind]
}

func ptrDecoderFunc(typ reflect.Type) decoderFunc {
	decoder := getDecoder(typ.Elem())
	return func(d *Decoder, v reflect.Value) error {
		if d.hasNilCode() {
			if !v.IsNil() {
				v.Set(reflect.Zero(v.Type()))
			}
			return d.DecodeNil()
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return decoder(d, v.Elem())
	}
}

var valueDecoders []decoderFunc

//nolint:gochecknoinits
func init() {
	valueDecoders = []decoderFunc{
		reflect.Bool:       decodeBoolValue,
		reflect.Int:        decodeIntValue,
		reflect.Int8:       decodeInt8Value,
		reflect.Int16:      decodeInt16Value,
		reflect.Int32:      decodeInt32Value,
		reflect.Int64:      decodeInt64Value,
		reflect.Uint:       decodeUintValue,
		reflect.Uint8:      decodeUint8Value,
		reflect.Uint16:     decodeUint16Value,
		reflect.Uint32:     decodeUint32Value,
		reflect.Uint64:     decodeUint64Value,
		reflect.Float32:    decodeFloat32Value,
		reflect.Float64:    decodeFloat64Value,
		reflect.Complex64:  decodeUnsupportedValue,
		reflect.Complex128: decodeUnsupportedValue,
		// reflect.Array:         decodeArrayValue,
		reflect.Chan: decodeUnsupportedValue,
		reflect.Func: decodeUnsupportedValue,
		//reflect.Interface:     decodeInterfaceValue,
		reflect.Map:           decodeMapValue,
		reflect.Ptr:           decodeUnsupportedValue,
		reflect.Slice:         decodeSliceValue,
		reflect.String:        decodeStringValue,
		reflect.Struct:        decodeStructValue,
		reflect.UnsafePointer: decodeUnsupportedValue,
	}
}

func decodeStringValue(d *Decoder, v reflect.Value) error {
	s, err := d.DecodeString()
	if err != nil {
		return err
	}
	v.SetString(s)
	return nil
}

func decodeUnsupportedValue(d *Decoder, v reflect.Value) error {
	return fmt.Errorf("phpserialize: Decode(unsupported %s)", v.Type())
}

func decodeBoolValue(d *Decoder, v reflect.Value) error {
	flag, err := d.DecodeBool()
	if err != nil {
		return err
	}
	v.SetBool(flag)
	return nil
}

func decodeIntValue(d *Decoder, v reflect.Value) error {
	return decodeSignedIntValue(d, v, bits.UintSize)
}

func decodeInt8Value(d *Decoder, v reflect.Value) error {
	return decodeSignedIntValue(d, v, 8)
}

func decodeInt16Value(d *Decoder, v reflect.Value) error {
	return decodeSignedIntValue(d, v, 16)
}

func decodeInt32Value(d *Decoder, v reflect.Value) error {
	return decodeSignedIntValue(d, v, 32)
}

func decodeInt64Value(d *Decoder, v reflect.Value) error {
	return decodeSignedIntValue(d, v, 64)
}

func decodeSignedIntValue(d *Decoder, v reflect.Value, bitSize int) error {
	n, err := d.DecodeSignedInt(bitSize)
	if err != nil {
		// TODO: should this wrapErr so the prefix phpserialize is maintained?
		return err
	}
	v.SetInt(n)
	return nil
}

func decodeFloat32Value(d *Decoder, v reflect.Value) error {
	n, err := d.DecodeFloat(32)
	if err != nil {
		// TODO: should this wrapErr so the prefix phpserialize is maintained?
		return err
	}
	v.SetFloat(n)
	return nil
}

func decodeFloat64Value(d *Decoder, v reflect.Value) error {
	n, err := d.DecodeFloat64()
	if err != nil {
		// TODO: should this wrapErr so the prefix phpserialize is maintained?
		return err
	}
	v.SetFloat(n)
	return nil
}

func decodeUintValue(d *Decoder, v reflect.Value) error {
	return decodeUnsignedIntValue(d, v, bits.UintSize)
}

func decodeUint8Value(d *Decoder, v reflect.Value) error {
	return decodeUnsignedIntValue(d, v, 8)
}

func decodeUint16Value(d *Decoder, v reflect.Value) error {
	return decodeUnsignedIntValue(d, v, 16)
}

func decodeUint32Value(d *Decoder, v reflect.Value) error {
	return decodeUnsignedIntValue(d, v, 32)
}

func decodeUint64Value(d *Decoder, v reflect.Value) error {
	return decodeUnsignedIntValue(d, v, 64)
}

func decodeUnsignedIntValue(d *Decoder, v reflect.Value, bitSize int) error {
	n, err := d.DecodeUnsignedInt(bitSize)
	if err != nil {
		// TODO: should this wrapErr so the prefix phpserialize is maintained?
		return err
	}
	v.SetUint(n)
	return nil
}
