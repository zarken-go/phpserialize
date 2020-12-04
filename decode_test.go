package phpserialize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"strings"
	"testing"
)

func TestUnmarshalInt8(t *testing.T) {
	var v int8
	container := struct {
		Value int8 `php:"v"`
	}{}

	assert.Nil(t, Unmarshal([]byte(`i:65;`), &v))
	assert.Equal(t, int8(65), v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";i:123;}`), &container))
	assert.Equal(t, int8(123), container.Value)

	assert.EqualError(t, Unmarshal([]byte(`i:128;`), &v), `strconv.ParseInt: parsing "128": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";i:-129;}`), &container), `strconv.ParseInt: parsing "-129": value out of range`)
}

func TestUnmarshalInt16(t *testing.T) {
	var v int16
	container := struct {
		Value int16 `php:"v"`
	}{}

	assert.Nil(t, Unmarshal([]byte(`i:65;`), &v))
	assert.Equal(t, int16(65), v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";i:123;}`), &container))
	assert.Equal(t, int16(123), container.Value)

	assert.EqualError(t, Unmarshal([]byte(`i:32768;`), &v), `strconv.ParseInt: parsing "32768": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";i:-32769;}`), &container), `strconv.ParseInt: parsing "-32769": value out of range`)
}

func TestUnmarshalInt32(t *testing.T) {
	var v int32
	container := struct {
		Value int32 `php:"v"`
	}{}

	assert.Nil(t, Unmarshal([]byte(`i:65;`), &v))
	assert.Equal(t, int32(65), v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";i:123;}`), &container))
	assert.Equal(t, int32(123), container.Value)

	assert.EqualError(t, Unmarshal([]byte(`i:2147483648;`), &v), `strconv.ParseInt: parsing "2147483648": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";i:-2147483649;}`), &container), `strconv.ParseInt: parsing "-2147483649": value out of range`)
}

func TestUnmarshalInt64(t *testing.T) {
	var v int64
	container := struct {
		Value int64 `php:"v"`
	}{}

	assert.Nil(t, UnmarshalString(fmt.Sprintf(`i:%d;`, math.MinInt64), &v))
	assert.Equal(t, int64(math.MinInt64), v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";i:123;}`), &container))
	assert.Equal(t, int64(123), container.Value)

	assert.EqualError(t, Unmarshal([]byte(`i:9223372036854775808;`), &v), `strconv.ParseInt: parsing "9223372036854775808": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";i:-9223372036854775809;}`), &container), `strconv.ParseInt: parsing "-9223372036854775809": value out of range`)
}

// assume testing happens on 64-bit system
func TestUnmarshalInt(t *testing.T) {
	var v int
	container := struct {
		Value int `php:"v"`
	}{}

	assert.Nil(t, UnmarshalString(fmt.Sprintf(`i:%d;`, math.MinInt64), &v))
	assert.Equal(t, math.MinInt64, v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";i:123;}`), &container))
	assert.Equal(t, 123, container.Value)

	assert.EqualError(t, Unmarshal([]byte(`i:9223372036854775808;`), &v), `strconv.ParseInt: parsing "9223372036854775808": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";i:-9223372036854775809;}`), &container), `strconv.ParseInt: parsing "-9223372036854775809": value out of range`)
}

func TestUnmarshalBool(t *testing.T) {
	var v bool
	container := struct {
		Value bool `php:"v"`
	}{}

	assert.Nil(t, UnmarshalString(`b:1;`, &v))
	assert.True(t, v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";b:1;}`), &container))
	assert.True(t, container.Value)

	assert.EqualError(t, Unmarshal([]byte(`b:2;`), &v), `phpserialize: Decode(invalid boolean value)`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";b:2;}`), &container), `phpserialize: Decode(invalid boolean value)`)
}

func TestUnmarshalFloat32(t *testing.T) {
	var v float32
	container := struct {
		Value float32 `php:"v"`
	}{}

	assert.Nil(t, UnmarshalString(`d:15.285325;`, &v))
	assert.Equal(t, float32(15.285325), v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";d:15235.12825;}`), &container))
	assert.Equal(t, float32(15235.12825), container.Value)

	assert.EqualError(t, Unmarshal([]byte(`d:3.402823466e+50;`), &v), `strconv.ParseFloat: parsing "3.402823466e+50": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";d:3.402823466e+50;}`), &container), `strconv.ParseFloat: parsing "3.402823466e+50": value out of range`)
}

func TestUnmarshalFloat64(t *testing.T) {
	var v float64
	container := struct {
		Value float64 `php:"v"`
	}{}

	assert.Nil(t, UnmarshalString(`d:15.285325;`, &v))
	assert.Equal(t, 15.285325, v)

	assert.Nil(t, Unmarshal([]byte(`a:1:{s:1:"v";d:15235.12825;}`), &container))
	assert.Equal(t, 15235.12825, container.Value)

	assert.Nil(t, UnmarshalString(`d:INF;`, &v))
	assert.Equal(t, math.Inf(1), v)

	assert.Nil(t, UnmarshalString(`d:-INF;`, &v))
	assert.Equal(t, math.Inf(-1), v)

	assert.Nil(t, UnmarshalString(`d:NAN;`, &v))
	assert.True(t, math.IsNaN(v))

	assert.EqualError(t, Unmarshal([]byte(`d:3.402823466e+325;`), &v), `strconv.ParseFloat: parsing "3.402823466e+325": value out of range`)
	assert.EqualError(t, Unmarshal([]byte(`a:1:{s:1:"v";d:3.402823466e+325;}`), &container), `strconv.ParseFloat: parsing "3.402823466e+325": value out of range`)
}

func TestUnmarshalSlice(t *testing.T) {
	var i []int
	assert.Nil(t, UnmarshalString(`a:3:{i:0;i:1;i:1;i:3;i:2;i:5;}`, &i))

	// TODO: special case
	// var s []string
	// assert.Nil(t, UnmarshalString(`a:3:{i:0;s:3:"one";i:1;s:5:"three";i:2;s:4:"five";}`, &s))
}

func TestUnmarshalSliceOfMaps(t *testing.T) {
	var m []map[string]string
	assert.Nil(t, UnmarshalString(`a:2:{i:0;a:2:{s:2:"id";s:1:"1";s:5:"value";s:3:"One";}i:1;a:2:{s:2:"id";s:1:"2";s:5:"value";s:3:"Two";}}`, &m))
	assert.Len(t, m, 2)
	assert.Equal(t, `1`, m[0][`id`])
	assert.Equal(t, `One`, m[0][`value`])
	assert.Equal(t, `2`, m[1][`id`])
	assert.Equal(t, `Two`, m[1][`value`])
}

func TestDecoder_DecodeFloat(t *testing.T) {
	d := NewDecoder(strings.NewReader(`b:1;`))
	v, err := d.DecodeFloat(64)
	assert.Zero(t, v)
	assert.EqualError(t, err, `phpserialize: Decode(expected byte 'd' found 'b')`)

	d = NewDecoder(strings.NewReader(`d:`))
	v, err = d.DecodeFloat(64)
	assert.Zero(t, v)
	assert.Equal(t, io.EOF, err)
}
