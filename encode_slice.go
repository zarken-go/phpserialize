package phpserialize

import "reflect"

func encodeSliceValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}
	return encodeArrayValue(e, v)
}

func encodeArrayValue(e *Encoder, v reflect.Value) error {
	l := v.Len()
	if err := e.writeArrayPrefixLen(l); err != nil {
		return err
	}
	for i := 0; i < l; i++ {
		if err := e.EncodeInt64(int64(i)); err != nil {
			return err
		}
		if err := e.EncodeValue(v.Index(i)); err != nil {
			return err
		}
	}
	return e.writeBytes('}')
}
