package phpserialize

import "reflect"

func encodeMapValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.EncodeNil()
	}

	if err := e.writeArrayPrefixLen(v.Len()); err != nil {
		return err
	}

	for _, key := range v.MapKeys() {
		if err := e.EncodeValue(key); err != nil {
			return err
		}
		if err := e.EncodeValue(v.MapIndex(key)); err != nil {
			return err
		}
	}

	return e.writeBytes('}')
}

func (e *Encoder) writeArrayPrefixLen(len int) error {
	if err := e.writeBytes('a', ':'); err != nil {
		return err
	}
	if err := e.writeInt(len); err != nil {
		return err
	}
	return e.writeBytes(':', '{')
}

func encodeStructValue(e *Encoder, strct reflect.Value) error {
	structFields := structs.Fields(strct.Type(), `php`) // e.structTag)
	/*if e.flags&arrayEncodedStructsFlag != 0 || structFields.AsArray {
		return encodeStructValueAsArray(e, strct, structFields.List)
	}*/
	fields := structFields.OmitEmpty(strct)

	if err := e.writeArrayPrefixLen(len(fields)); err != nil {
		return err
	}

	for _, f := range fields {
		if err := e.EncodeString(f.name); err != nil {
			return err
		}
		if err := f.EncodeValue(e, strct); err != nil {
			return err
		}
	}

	return e.writeBytes('}')
}
