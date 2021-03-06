package phpserialize

import (
	"fmt"
	"reflect"
)

func decodeStructValue(d *Decoder, v reflect.Value) error {
	arrayLen, err := d.decodeArrayLen()
	if err != nil {
		return err
	}

	fields := structs.Fields(v.Type(), defaultStructTag)
	for i := 0; i < arrayLen; i++ {
		name, err := d.DecodeString()
		if err != nil {
			return err
		}

		if f := fields.Map[name]; f != nil {
			if err := f.DecodeValue(d, v); err != nil {
				return err
			}
			} else if d.flags&disallowUnknownFieldsFlag != 0 {
				return fmt.Errorf("phpserialize: unknown field %q", name)
			//} else if err := d.Skip(); err != nil {
			//	return err
		}
	}

	return d.skipExpected('}')
}
