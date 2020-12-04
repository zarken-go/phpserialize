package phpserialize

import "reflect"

func decodeMapValue(d *Decoder, v reflect.Value) error {
	n, err := d.decodeArrayLen()
	if err != nil {
		return err
	}

	typ := v.Type()
	if n == -1 {
		v.Set(reflect.Zero(typ))
		return nil
	}

	if v.IsNil() {
		v.Set(reflect.MakeMap(typ))
	}
	if n == 0 {
		return nil
	}

	if err := d.decodeTypedMapValue(v, n); err != nil {
		return err
	}

	return d.skipExpected('}')
}

func (d *Decoder) decodeTypedMapValue(v reflect.Value, n int) error {
	typ := v.Type()
	keyType := typ.Key()
	valueType := typ.Elem()

	for i := 0; i < n; i++ {
		mk := reflect.New(keyType).Elem()
		if err := d.DecodeValue(mk); err != nil {
			return err
		}

		mv := reflect.New(valueType).Elem()
		if err := d.DecodeValue(mv); err != nil {
			return err
		}

		v.SetMapIndex(mk, mv)
	}

	return nil
}
