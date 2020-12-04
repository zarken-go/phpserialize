package phpserialize

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"strconv"
)

type writer interface {
	io.Writer
	WriteByte(byte) error
}

type byteWriter struct {
	io.Writer
}

func newByteWriter(w io.Writer) byteWriter {
	return byteWriter{
		Writer: w,
	}
}

func (bw byteWriter) WriteByte(c byte) error {
	_, err := bw.Write([]byte{c})
	return err
}

// Marshal returns the MessagePack encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	//enc := GetEncoder()
	enc := Encoder{}

	var buf bytes.Buffer
	//enc.Reset(&buf)
	enc.resetWriter(&buf)

	err := enc.Encode(v)
	b := buf.Bytes()

	// PutEncoder(enc)

	if err != nil {
		return nil, err
	}
	return b, err
}

type Encoder struct {
	w writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{}
	e.resetWriter(w)
	return e
}

func (e *Encoder) resetWriter(w io.Writer) {
	if bw, ok := w.(writer); ok {
		e.w = bw
	} else {
		e.w = newByteWriter(w)
	}
}

func (e *Encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return e.EncodeNil()
	case string:
		return e.EncodeString(v)
	case []byte:
		return e.EncodeBytes(v)
	case int:
		return e.EncodeInt64(int64(v))
	case int8:
		return e.EncodeInt64(int64(v))
	case int16:
		return e.EncodeInt64(int64(v))
	case int32:
		return e.EncodeInt64(int64(v))
	case int64:
		return e.EncodeInt64(v)
	/*case uint:
		return e.EncodeUint(uint64(v))
	case uint64:
		return e.encodeUint64Cond(v)*/
	case bool:
		return e.EncodeBool(v)
	case float32:
		return e.EncodeFloat64(float64(v))
	case float64:
		return e.EncodeFloat64(v)
		/*case time.Duration:
			return e.encodeInt64Cond(int64(v))
		case time.Time:
			return e.EncodeTime(v)
		*/
	}

	return e.EncodeValue(reflect.ValueOf(v))
}

func (e *Encoder) EncodeValue(v reflect.Value) error {
	fn := getEncoder(v.Type())
	return fn(e, v)
}

func (e *Encoder) EncodeNil() error {
	return e.writeBytes('N', ';')
}

func (e *Encoder) EncodeString(v string) error {
	return e.EncodeBytes([]byte(v))
}

func (e *Encoder) EncodeBytes(v []byte) error {
	if err := e.writeBytes('s', ':'); err != nil {
		return err
	}
	if err := e.writeInt(len(v)); err != nil {
		return err
	}
	if err := e.writeBytes(':', '"'); err != nil {
		return err
	}
	if err := e.write(v); err != nil {
		return err
	}
	return e.writeBytes('"', ';')
}

func (e *Encoder) EncodeInt64(v int64) error {
	if err := e.writeBytes('i', ':'); err != nil {
		return err
	}
	if err := e.writeInt64(v); err != nil {
		return err
	}
	return e.writeBytes(';')
}

func (e *Encoder) EncodeBool(v bool) error {
	if err := e.writeBytes('b', ':'); err != nil {
		return err
	}
	var boolChar byte = '0'
	if v {
		boolChar = '1'
	}
	if err := e.writeBytes(boolChar); err != nil {
		return err
	}
	return e.writeBytes(';')
}

func (e *Encoder) EncodeUint64(v uint64) error {
	if err := e.writeBytes('i', ':'); err != nil {
		return err
	}
	if err := e.writeUint64(v); err != nil {
		return err
	}
	return e.writeBytes(';')
}

func (e *Encoder) EncodeFloat64(v float64) error {
	if err := e.writeBytes('d', ':'); err != nil {
		return err
	}
	if math.IsInf(v, -1) {
		if err := e.writeString(`-INF`); err != nil {
			return err
		}
	} else if math.IsInf(v, 1) {
		if err := e.writeString(`INF`); err != nil {
			return err
		}
	} else if math.IsNaN(v) {
		if err := e.writeString(`NAN`); err != nil {
			return err
		}
	} else {
		if err := e.writeFloat64(v); err != nil {
			return err
		}
	}

	return e.writeBytes(';')
}

func (e *Encoder) write(b []byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) writeBytes(b ...byte) error {
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) writeString(s string) error {
	// _, err := e.w.Write(stringToBytes(s))
	_, err := e.w.Write([]byte(s))
	return err
}

func (e *Encoder) writeInt(v int) error {
	return e.writeString(strconv.Itoa(v))
}

func (e *Encoder) writeInt64(v int64) error {
	return e.writeString(strconv.FormatInt(v, 10))
}

func (e *Encoder) writeUint64(v uint64) error {
	return e.writeString(strconv.FormatUint(v, 10))
}

func (e *Encoder) writeFloat64(v float64) error {
	return e.writeString(strconv.FormatFloat(v, 'f', -1, 64))
}
