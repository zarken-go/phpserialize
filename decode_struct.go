package phpserialize

import "reflect"

func decodeStructValue(d *Decoder, v reflect.Value) error {
	if err := d.skipExpected('a', ':'); err != nil {
		return err
	}
	arrayLen, err := d.readUntilLen()
	if err != nil {
		return err
	}

	if err := d.skipExpected('{'); err != nil {
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
			/*} else if d.flags&disallowUnknownFieldsFlag != 0 {
				return fmt.Errorf("phpserialize: unknown field %q", name)
			} else if err := d.Skip(); err != nil {
				return err*/
		}
	}

	if err := d.skipExpected('}'); err != nil {
		return err
	}

	return nil
}
