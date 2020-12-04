package phpserialize

import (
	"fmt"
	"reflect"
)

const (
	sliceAllocLimit = 1e4
)

func decodeSliceValue(d *Decoder, v reflect.Value) error {
	n, err := d.decodeArrayLen()
	if err != nil {
		return err
	}

	if n == -1 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	if n == 0 && v.IsNil() {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		return nil
	}

	if v.Cap() >= n {
		v.Set(v.Slice(0, n))
	} else if v.Len() < v.Cap() {
		v.Set(v.Slice(0, v.Cap()))
	}

	for i := 0; i < n; i++ {
		if i >= v.Len() {
			v.Set(growSliceValue(v, n))
		}
		elem := v.Index(i)
		decodedIndex, err := d.DecodeInt()
		if err != nil {
			return err
		}
		if decodedIndex != i {
			return fmt.Errorf(`phpserialize: Decode(expected offset '%d' found '%d')`, i, decodedIndex)
		}
		if err := d.DecodeValue(elem); err != nil {
			return err
		}
	}

	return nil
}

func growSliceValue(v reflect.Value, n int) reflect.Value {
	diff := n - v.Len()
	if diff > sliceAllocLimit {
		diff = sliceAllocLimit
	}
	v = reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), diff, diff))
	return v
}
