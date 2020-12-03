package phpserialize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
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
