package phpserialize

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Decoder struct {
	s     io.ByteScanner
	flags uint32
}

const (
	disallowUnknownFieldsFlag uint32 = 1 << iota
)

const (
	// bytesAllocLimit = 1e6 // 1mb
	// sliceAllocLimit = 1e4
	maxMapSize = 1e6
)

var (
	ErrUnsupported = errors.New(`unsupported target type`)
)

type bufReader interface {
	// io.Reader
	io.ByteScanner
}

func Unmarshal(data []byte, v interface{}) error {
	d := NewDecoder(bytes.NewReader(data))
	return d.Decode(v)
}

func UnmarshalString(data string, v interface{}) error {
	d := NewDecoder(strings.NewReader(data))
	return d.Decode(v)
}

func NewDecoder(r io.Reader) *Decoder {
	d := new(Decoder)
	d.resetReader(r)
	return d
}

func (d *Decoder) resetReader(r io.Reader) {
	if br, ok := r.(bufReader); ok {
		//d.r = br
		d.s = br
	} else {
		br := bufio.NewReader(r)
		//d.r = br
		d.s = br
	}
}

//nolint:gocyclo
func (d *Decoder) Decode(v interface{}) error {
	var err error
	switch v := v.(type) {
	case *string:
		if v != nil {
			*v, err = d.DecodeString()
			return err
		}
	case *[]byte:
		if v != nil {
			return ErrUnsupported // d.decodeBytesPtr(v)
		}
	case *int:
		if v != nil {
			*v, err = d.DecodeInt()
			return err
		}
	case *int8:
		if v != nil {
			*v, err = d.DecodeInt8()
			return err
		}
	case *int16:
		if v != nil {
			*v, err = d.DecodeInt16()
			return err
		}
	case *int32:
		if v != nil {
			*v, err = d.DecodeInt32()
			return err
		}
	case *int64:
		if v != nil {
			*v, err = d.DecodeInt64()
			return err
		}
	case *uint:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeUint()
			return err
		}*/
	case *uint8:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeUint8()
			return err
		}*/
	case *uint16:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeUint16()
			return err
		}*/
	case *uint32:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeUint32()
			return err
		}*/
	case *uint64:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeUint64()
			return err
		}*/
	case *bool:
		if v != nil {
			*v, err = d.DecodeBool()
			return err
		}
	case *float32:
		if v != nil {
			*v, err = d.DecodeFloat32()
			return err
		}
	case *float64:
		if v != nil {
			*v, err = d.DecodeFloat64()
			return err
		}
	case *[]string:
		return d.decodeStringSlicePtr(v)
	case *map[string]string:
		return d.decodeMapStringStringPtr(v)
	case *map[string]interface{}:
		return ErrUnsupported // d.decodeMapStringInterfacePtr(v)
	case *time.Duration:
		if v != nil {
			vv, err := d.DecodeInt64()
			*v = time.Duration(vv)
			return err
		}
	case *time.Time:
		return ErrUnsupported
		/*if v != nil {
			*v, err = d.DecodeTime()
			return err
		}*/
	}

	vv := reflect.ValueOf(v)
	if !vv.IsValid() {
		return errors.New("phpserialize: Decode(nil)")
	}
	if vv.Kind() != reflect.Ptr {
		return fmt.Errorf("phpserialize: Decode(non-pointer %T)", v)
	}
	if vv.IsNil() {
		return fmt.Errorf("phpserialize: Decode(non-settable %T)", v)
	}

	vv = vv.Elem()
	if vv.Kind() == reflect.Interface {
		if !vv.IsNil() {
			vv = vv.Elem()
			if vv.Kind() != reflect.Ptr {
				return fmt.Errorf("phpserialize: Decode(non-pointer %s)", vv.Type().String())
			}
		}
	}

	return d.DecodeValue(vv)
}

func (d *Decoder) PeekCode() (byte, error) {
	c, err := d.s.ReadByte()
	if err != nil {
		return 0, err
	}
	return c, d.s.UnreadByte()
}

func (d *Decoder) hasNilCode() bool {
	code, err := d.PeekCode()
	return err == nil && code == 'N'
}

func (d *Decoder) DecodeNil() error {
	return d.skipExpected('N', ';')
}

func (d *Decoder) DecodeValue(v reflect.Value) error {
	decode := getDecoder(v.Type())
	if decode == nil {
		return fmt.Errorf(`phpserialize: could not find decoder for: %s`, v.Type().String())
	}
	return decode(d, v)
}

/**
  b:1;
  b:0;
*/
func (d *Decoder) DecodeBool() (bool, error) {
	if err := d.skipExpected('b', ':'); err != nil {
		return false, err
	}
	v, err := d.s.ReadByte()
	if err != nil {
		return false, err
	}
	if err := d.skipExpected(';'); err != nil {
		return false, err
	}
	switch v {
	case '1':
		return true, nil
	case '0':
		return false, nil
	default:
		return false, errors.New(`phpserialize: Decode(invalid boolean value)`)
	}
}

/**
  i:685230;
  i:-685230;
*/

func (d *Decoder) DecodeInt() (int, error) {
	v, err := d.DecodeSignedInt(bits.UintSize)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (d *Decoder) DecodeInt8() (int8, error) {
	v, err := d.DecodeSignedInt(8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

func (d *Decoder) DecodeInt16() (int16, error) {
	v, err := d.DecodeSignedInt(16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (d *Decoder) DecodeInt32() (int32, error) {
	v, err := d.DecodeSignedInt(32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (d *Decoder) DecodeInt64() (int64, error) {
	return d.DecodeSignedInt(64)
}

func (d *Decoder) DecodeSignedInt(bitSize int) (int64, error) {
	if err := d.skipExpected('i', ':'); err != nil {
		return 0, err
	}

	acc, err := d.readUntil(';')
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(acc), 10, bitSize)
}

func (d *Decoder) DecodeUnsignedInt(bitSize int) (uint64, error) {
	if err := d.skipExpected('i', ':'); err != nil {
		return 0, err
	}

	acc, err := d.readUntil(';')
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(acc), 10, bitSize)
}

/**
d:685230.15;
d:INF;
d:-INF;
d:NAN;
*/
func (d *Decoder) DecodeFloat64() (float64, error) {
	return d.DecodeFloat(64)
}

func (d *Decoder) DecodeFloat32() (float32, error) {
	v, err := d.DecodeFloat(32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

func (d *Decoder) DecodeFloat(bitSize int) (float64, error) {
	if err := d.skipExpected('d', ':'); err != nil {
		return 0, err
	}

	acc, err := d.readUntil(';')
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(string(acc), bitSize)
}

func (d *Decoder) DecodeString() (string, error) {
	if err := d.skipExpected('s', ':'); err != nil {
		return ``, err
	}
	strLen, err := d.readUntilLen()
	if err != nil {
		return ``, err
	}
	if err := d.skipExpected('"'); err != nil {
		return ``, err
	}
	acc := make([]byte, strLen)
	for x := 0; x < strLen; x++ {
		b, err := d.s.ReadByte()
		if err != nil {
			return ``, err
		}
		acc[x] = b
	}
	if err := d.skipExpected('"', ';'); err != nil {
		return ``, err
	}

	return string(acc), nil
}

func (d *Decoder) readUntil(v byte) ([]byte, error) {
	var acc []byte
	for {
		b, err := d.s.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == v {
			break
		}
		acc = append(acc, b)
	}
	return acc, nil
}

func (d *Decoder) readUntilLen() (int, error) {
	acc, err := d.readUntil(':')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(acc))
}

func (d *Decoder) skipExpected(expected ...byte) error {
	for _, e := range expected {
		c, err := d.s.ReadByte()
		if err != nil {
			return err
		}
		if c != e {
			return fmt.Errorf(`phpserialize: Decode(expected byte '%c' found '%c')`, e, c)
		}
	}
	return nil
}

func (d *Decoder) decodeArrayLen() (int, error) {
	if err := d.skipExpected('a', ':'); err != nil {
		return 0, err
	}
	n, err := d.readUntilLen()
	if err != nil {
		return 0, err
	}
	if err := d.skipExpected('{'); err != nil {
		return 0, err
	}
	return n, nil
}

func min(a, b int) int { //nolint:unparam
	if a <= b {
		return a
	}
	return b
}
